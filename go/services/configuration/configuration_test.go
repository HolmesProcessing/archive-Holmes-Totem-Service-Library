package configuration

import (
    "io/ioutil"
    "os"
    "testing"
)

func TestConfiguration(t *testing.T) {
    file, err := ioutil.TempFile("", "holmes-totem-service-library-testing-go_Configuration")
    if err != nil {
        t.Error("Unable to create temporary file for testing")
        return
    }
    filepath := file.Name()
    defer os.Remove(filepath)

    file.Write([]byte(`
        {
            "key1": "value1",
            "section2": {
                "key1": "value1",
                "keY2": 1337,
                "subSeCtIon1": {
                    "KEY1": 10.5
                }
            }
        }
    `))
    file.Close()

    var testconfig struct {
        Key1     string
        Section2 struct {
            Key1        string
            Key2        int
            Subsection1 struct {
                Key1 float32
            }
        }
    }

    Parse(&testconfig, filepath)

    success :=
        (testconfig.Key1 == "value1") &&
        (testconfig.Section2.Key1 == "value1") &&
        (testconfig.Section2.Key2 == 1337) &&
        (testconfig.Section2.Subsection1.Key1 == 10.5)

    if !success {
        t.Error("Configuration test failed")
        return
    }
}
