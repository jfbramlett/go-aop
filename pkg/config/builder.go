package config

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"io/ioutil"
	"os"
	"strings"
)

type Source interface {
	Read() map[string]interface{}
}

func NewFileSource(filename string) Source {
	return &FileSource{filename: filename}
}

type FileSource struct {
	filename	string
}

func (f *FileSource) Read() map[string]interface{} {
	cfg := make(map[string]interface{})

	info, err := os.Stat(f.filename)
	if os.IsNotExist(err) || info.IsDir() {
		return cfg
	}

	fileContent, err := ioutil.ReadFile(f.filename)
	if err != nil {
		return cfg
	}
	err = json.Unmarshal(fileContent, cfg)

	return cfg
}

func NewEnvSource() Source {
	return &EnvSource{}
}

type EnvSource struct {
}

func (e *EnvSource) Read() map[string]interface{} {
	cfg := make(map[string]interface{})

	for _, e := range os.Environ() {
		kv := strings.Split(e, "=")
		if len(kv) != 2 {
			continue
		}
		cfg[kv[0]] = kv[1]
	}

	return cfg
}

func NewArgSource(flagSet *pflag.FlagSet) Source {
	return &ArgSource{flagSet: flagSet}
}

type ArgSource struct {
	flagSet 	*pflag.FlagSet
}

func (a *ArgSource) Read() map[string]interface{} {
	flagJson := ""
	a.flagSet.VisitAll(func(f *pflag.Flag) {
		if len(flagJson) > 0 {
			flagJson = fmt.Sprintf("%s,", flagJson)
		}
		if f.Value.Type() == "string" {
			flagJson = fmt.Sprintf("%s\"%s\": \"%s\"", flagJson, f.Name, f.Value.String())
		} else {
			flagJson = fmt.Sprintf("%s\"%s\": %s", flagJson, f.Name, f.Value.String())
		}
	})
	flagJson = fmt.Sprintf("{%s}", flagJson)

	cfg := make(map[string]interface{})
	_ = json.Unmarshal([]byte(flagJson), &cfg)

	return cfg
}

func NewBuilder() *Builder {
	return &Builder{sources: make([]Source, 0)}
}

type Builder struct {
	sources []Source
}

func (b *Builder) WithSource(s Source) *Builder {
	b.sources = append(b.sources, s)
	return b
}

func (b *Builder) Build(obj interface{}) error {
	cfg := make(map[string]interface{})
	for _, s := range b.sources {
		sCfg := s.Read()
		for k, v := range sCfg {
			cfg[k] = v
		}
	}

	jsonStr, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonStr, obj)
	if err != nil {
		return err
	}

	return nil
}