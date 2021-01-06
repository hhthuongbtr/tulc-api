package utils


import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
)

func GetMd5FromFile(path string) (md5string string, err error)  {
	file, err := os.Open(path)

	if err != nil {
		return "", err
	}

	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		return "", err
	}
	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]
	//Convert the bytes to a string
	returnMD5String := hex.EncodeToString(hashInBytes)

	return returnMD5String, nil
}

func GetFileSizeInByte(path string) (fileSize int64, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	// get the size
	size := fi.Size()
	return size, nil
}


type MyFile struct {
	Path	string
}

func (mf *MyFile) Read() (string, error) {
	data, err := ioutil.ReadFile(mf.Path)
	if err != nil {
		return "", err
	}
	return string(data), err
}

// Exists reports whether the named file or directory exists.
func (mf *MyFile) Exists() bool {
	if _, err := os.Stat(mf.Path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (mf *MyFile) GetMd5FromFile(path string) (md5string string, err error)  {
	file, err := os.Open(path)

	if err != nil {
		return "", err
	}

	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		return "", err
	}
	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]
	//Convert the bytes to a string
	returnMD5String := hex.EncodeToString(hashInBytes)

	return returnMD5String, nil
}

func (mf *MyFile) GetFileSizeInByte(path string) (fileSize int64, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	// get the size
	size := fi.Size()
	return size, nil
}
