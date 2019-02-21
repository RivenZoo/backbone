package main

import (
	"bytes"
	"github.com/RivenZoo/backbone/logger"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	parseFlagConfig()

	if config.Debug {
		logger.SetLogLevel(logger.DEBUG)
	} else {
		logger.SetLogLevel(logger.INFO)
	}

	logger.Debugf("config %s", config)

	var err error
	src := config.LocalPath
	dst := config.OutputDir

	if config.LocalPath == "" {
		src, err = fetchRepo(config.GitRepoURL, config.GitVersion)
		if err != nil {
			logger.Errorf("fetch repo %s version %s error %v", config.GitRepoURL, config.GitVersion,
				err)
			os.Exit(-1)
		}
		defer func() {
			os.RemoveAll(src)
		}()
	}

	err = copyDirRecursively(src, dst)
	if err != nil {
		logger.Errorf("copy from %s to %s error %v", src, dst, err)
		os.Exit(-1)
	}

	execOut, err := execProgram(config.ExecAfterCreate, dst)
	if err != nil {
		logger.Errorf("exec %s error %v, skipped", config.ExecAfterCreate, err)
	} else {
		logger.Infof("exec output\n %s", string(execOut))
	}
	writeConfig(dst)
}

func copyDirRecursively(srcDir, outputDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Errorf("handle file %s error %v", path, err)
			// end walk
			return err
		}

		if info.IsDir() {
			err = copyDirPath(srcDir, outputDir, path)
		} else {
			err = copyFile(srcDir, outputDir, path)
		}
		return err
	})
}

func replacePathPrefix(srcDir, dstDir, path string) string {
	pathSeg := make([]string, 0)
	srcDir = filepath.Clean(srcDir)
	path = filepath.Clean(path)
	if path == srcDir {
		return filepath.Clean(dstDir)
	}
	for {
		f := filepath.Base(path)
		pathSeg = append(pathSeg, f)
		d := filepath.Dir(path)
		if d == srcDir || d == "." || d == "/" { // d == "." or "/"
			break
		}
		path = d
	}
	pathSeg = append(pathSeg, dstDir)

	for i := len(pathSeg)/2 - 1; i >= 0; i-- {
		opp := len(pathSeg) - 1 - i
		pathSeg[i], pathSeg[opp] = pathSeg[opp], pathSeg[i]
	}
	return filepath.Join(pathSeg...)
}

func replacePathByTmplVar(path string) (string, error) {
	if strings.Index(path, "{{") == -1 {
		// no template variable
		return path, nil
	}
	t, err := template.New("dirTmpl").Parse(path)
	if err != nil {
		return path, err
	}
	buf := bytes.NewBuffer(make([]byte, 0))
	err = t.Execute(buf, config.TmplArgs)
	return buf.String(), err
}

func replaceFileByTmplVar(filePath string) (io.Reader, error) {
	if strings.ToLower(filepath.Ext(filePath)) != ".tmpl" {
		return os.Open(filePath)
	}
	t, err := template.New("fileTmpl").Parse(filePath)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(make([]byte, 0))
	err = t.Execute(buf, config.TmplArgs)
	return buf, err
}

func copyDirPath(srcDir, outputDir, dirPath string) error {
	outputPath := replacePathPrefix(srcDir, outputDir, dirPath)
	d, err := replacePathByTmplVar(outputPath)
	if err != nil {
		return err
	}
	return os.MkdirAll(d, 0755)
}

func detectFilePerm(filePath string) os.FileMode {
	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".sh", ".bash", ".zsh", ".py":
		return 0744
	default:
	}
	return 0644
}

func copyFile(srcDir, outputDir, filePath string) error {
	outputFile := replacePathPrefix(srcDir, outputDir, filePath)
	r, err := replaceFileByTmplVar(filePath)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outputFile, data, detectFilePerm(outputFile))
}
