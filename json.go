package gonfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type JsonConfig struct {
	Configurable
	Path string
}

func unmarshalJSONSection(output map[string]string, out map[string]interface{}, path string) interface{} {

	for k, v := range out {
		kPath := k
		if path != "" {
			kPath = fmt.Sprintf("%s:%s", path, k)
		}
		switch v := v.(type) {
		case int:
			output[kPath] = fmt.Sprintf("%d", v)
		case float64:
			output[kPath] = fmt.Sprintf("%f", v)
		case string:
			output[kPath] = v
		case bool:
			if v {
				output[kPath] = "true"
			} else {
				output[kPath] = "false"
			}
		case []interface{}:
			// What is the expected behavior on lists??
		case map[string]interface{}:
			unmarshalJSONSection(output, v, kPath)
		default:
			// weird
		}
	}
	return output

}

func unmarshalJson(bytes []byte) (map[string]string, error) {
	out := make(map[string]interface{})
	if err := json.Unmarshal(bytes, &out); err != nil {
		return nil, err
	}
	output := make(map[string]string)
	unmarshalJSONSection(output, out, "")
	return output, nil
}

// Returns a new WritableConfig backed by a json file at path.
// The file does not need to exist, if it does not exist the first Save call will create it.
func NewJsonConfig(path string, cfg ...Configurable) WritableConfig {
	if len(cfg) == 0 {
		cfg = append(cfg, NewMemoryConfig())
	}
	LoadConfig(cfg[0])
	conf := &JsonConfig{cfg[0], path}
	LoadConfig(conf)
	return conf
}

// Attempts to load the json configuration at JsonConfig.Path
// and Set them into the underlaying Configurable
func (self *JsonConfig) Load() (err error) {
	var data []byte = make([]byte, 1024)
	if data, err = ioutil.ReadFile(self.Path); err != nil {
		return err
	}
	out, err := unmarshalJson(data)
	if err != nil {
		return err
	}

	self.Configurable.Reset(out)
	return nil
}

// Attempts to save the configuration from the underlaying Configurable to json file at JsonConfig.Path
func (self *JsonConfig) Save() (err error) {
	b, err := json.Marshal(self.Configurable.All())
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(self.Path, b, 0600); err != nil {
		return err
	}

	return nil
}
