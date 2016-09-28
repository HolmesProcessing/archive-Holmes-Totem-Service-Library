package configuration

import (
    "encoding/json"
    "os"
    "bufio"
)

func check(e error, prefix string) {
    if e != nil {
        panic(prefix + e.Error())
    }
}

/*
 * Fill a struct given by a pointer via #dest with the values found at #path.
 * The file saved in #path needs to contain exactly one valid JSON struct.
 *
 * TODO: Set a field to be required by adding the tag maybe?
 *
 *      `jsonconfig:"required"`
 *
 */
func Parse(config interface{}, path string) {
    // Read the JSON file:
    file, err := os.Open(path);
    check(err, "Failed to open the config file: ")

    // Parse the JSON file:
    reader := bufio.NewReader(file)
    decoder := json.NewDecoder(reader)
    err = decoder.Decode(config)
    check(err, "Failed to parse the config file: ")
}
