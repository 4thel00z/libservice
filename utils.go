package libservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/monzo/typhon"
	"gopkg.in/dealancer/validate.v2"
)

const (
	multipartTag = "multipart"
)

func Default404Handler(app App) typhon.Service {
	return func(req typhon.Request) typhon.Response {
		// TODO: Change this body to a default 404 page
		response := req.Response(nil)
		response.StatusCode = 404
		return response
	}
}

func GenerateJSONValidator(i interface{}) *Validator {
	t := reflect.TypeOf(i)
	toValidate := reflect.New(t).Interface()

	validator := func(r typhon.Request) (interface{}, error) {

		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		// As if nothing has ever happened .. ( ͡° ͜ʖ ͡°)
		r.Body = ioutil.NopCloser(bytes.NewReader(content))

		err = json.Unmarshal(content, &toValidate)

		if err != nil {
			return nil, err
		}

		err = validate.Validate(toValidate)

		if err != nil {
			return nil, err
		}

		return toValidate, nil
	}

	return (*Validator)(&validator)
}

func getMultipartTags(t reflect.Type) map[string]bool {
	required := map[string]bool{}
	// Iterate over struct fields
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		key := field.Tag.Get("mapstructure")
		if key == "" {
			key = strcase.ToSnake(field.Name)
		}
		val := strings.Trim(strings.ToLower(field.Tag.Get(multipartTag)), " ")
		required[key] = val == "required"
	}
	return required
}
func GenerateMultipartValidator(i interface{}, maxMemoryBytes int64) *Validator {
	t := reflect.TypeOf(i)
	toValidate := reflect.New(t).Interface()

	validator := func(r typhon.Request) (interface{}, error) {
		err := r.ParseMultipartForm(maxMemoryBytes)

		if err != nil {
			return nil, err
		}
		tags := getMultipartTags(t)

		contents := map[string][]byte{}

		for tag, required := range tags {

			file, _, err := r.FormFile(tag)
			if err != nil && required {
				return nil, errors.New(fmt.Sprintf("you need to set %s", tag))
			}
			defer file.Close()
			content, err := ioutil.ReadAll(file)
			if err != nil && required {
				return nil, err
			}
			contents[tag] = content
		}
		err = mapstructure.Decode(contents, &toValidate)
		if err != nil {
			return nil, err
		}
		err = validate.Validate(toValidate)

		if err != nil {
			return nil, err
		}

		return toValidate, nil
	}

	return (*Validator)(&validator)
}

func GetCurrentDir() (dirAbsPath string, err error) {

	ex, err := os.Executable()

	if err == nil {
		dirAbsPath = filepath.Dir(ex)
		return dirAbsPath, err
	}

	exReal, err := filepath.EvalSymlinks(ex)

	if err != nil {
		return "", err
	}

	dirAbsPath = filepath.Dir(exReal)
	return dirAbsPath, err
}

func GetPackagePath(i interface{}) string {
	if i == nil {
		return ""
	}
	val := reflect.ValueOf(i)
	if val.Kind() == reflect.Ptr {
		return val.Elem().Type().PkgPath()
	}
	return val.Type().PkgPath()
}

func GetGoEnv() (map[string]string, error) {
	tool, err := exec.LookPath("go")
	if err != nil {
		return nil, err
	}

	out, err := exec.Command(tool, "env", "-json").CombinedOutput()

	if err != nil {
		return nil, err
	}

	var result map[string]string
	err = json.Unmarshal(out, &result)

	if err != nil {
		return nil, err
	}

	return result, err
}
