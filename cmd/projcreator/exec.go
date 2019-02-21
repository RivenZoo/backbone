package main

import (
	"os"
	"os/exec"
)

func parseExecCommandFields(execConfig string) []string {
	fields := make([]string, 0)
	startQuotation := false
	startSingleQuotation := false
	field := make([]byte, 0)
	for i := range execConfig {
		c := execConfig[i]
		switch c {
		case ' ', '\t', '\n':
			if startQuotation || startSingleQuotation {
				field = append(field, c)
			} else if len(field) > 0 {
				fields = append(fields, string(field))
				field = make([]byte, 0)
			}
		case '"':
			if startSingleQuotation {
				// no op
			} else if startQuotation {
				if len(field) > 0 {
					fields = append(fields, string(field))
					field = make([]byte, 0)
				}
				startQuotation = false
			} else {
				startQuotation = true
			}
		case '\'':
			if startQuotation {
				// no op
			} else if startSingleQuotation {
				if len(field) > 0 {
					fields = append(fields, string(field))
					field = make([]byte, 0)
				}
				startSingleQuotation = false
			} else {
				startSingleQuotation = true
			}
		default:
			field = append(field, c)
		}
	}
	if len(field) > 0 {
		fields = append(fields, string(field))
	}
	return fields
}

// execProgram if no output, return nil, nil
func execProgram(execCmd, workDir string) ([]byte, error) {
	if execCmd == "" {
		return nil, nil
	}
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	defer func() {
		os.Chdir(pwd)
	}()

	err = os.Chdir(workDir)
	if err != nil {
		return nil, err
	}

	cmds := parseExecCommandFields(execCmd)
	cmd := exec.Command(cmds[0], cmds[1:]...)
	return cmd.Output()
}
