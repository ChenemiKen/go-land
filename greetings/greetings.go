package greetings

import "fmt"
import "errors"
import "math/rand"


func Hello(name string) (string, error){
	if name == ""{
		return "", errors.New("empty name")
	}
	message := fmt.Sprintf(randonFormat(), name)
	return message, nil
}

func Hellos(names []string) (map[string]string, error){
	messages := make(map[string]string)
	for _, name := range names{
		message, err := Hello(name)

		if err != nil {
			return nil, err
		}

		messages[name] = message
	}
	return messages, nil
}

func randonFormat()string{
	formats := []string{
		"Hi, %v, Welcome!",
		"Great to see you, %v!",
		"Hail, %v!, well met!",
	} 

	return formats[rand.Intn(len(formats))]
}