package parser

import (
	"bufio"
	"os"
	"strings"

	"bitbucket.org/vservices/dark/logger"
)

var log = logger.New()

//LinesFromFile loads a file into an array of lines
//skipping empty lines and comments starting with "--" up to end of that line
//and trimming white space before/after the line text
//comments are published with the next significant line
//only comments at the end of the file will be written to a line without text
func LinesFromFile(filename string) (ILines, error) {
	inFile, err := os.Open(filename)
	if err != nil {
		return nil, log.Wrapf(err, "Cannot read file %s", filename)
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	lineNr := 0
	lines := NewLines()

	comment := ""
	for scanner.Scan() {
		lineNr++

		//trim spaces from both sides of the line
		lineText := strings.Trim(scanner.Text(), " \t")

		//trim comment starting anywhere in the line up to end of the line
		//	""					-> [""]					-> line=""
		//  "--comment"			-> ["", "comment"]		-> line=""
		//  "text"				-> ["text"]				-> line="text"
		//  "text --comment"	-> ["text", "comment"]	-> line="text"
		parts := strings.SplitN(lineText, "--", 2)
		if len(parts) > 0 {
			lineText = parts[0]
			if len(parts) > 1 {
				comment += parts[1]
			}
		}

		lineLen := len(lineText)
		if lineLen == 0 {
			//skip empty line
			continue
		}

		lines = lines.Append(lineNr, lineText, comment)
		comment = ""
	} //for each line

	if len(comment) > 0 {
		//comment at end of file
		lines = lines.Append(lineNr, "", comment)
	}

	return lines, nil
} //LinesFromFile()
