package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Conversation struct {
	Messages []Message `json:"messages"`
}

func readFilePath() string {
	fmt.Println("\n请拖入一个txt文件到命令行中，或输入文件路径后按回车：")
	fmt.Println("输入'exit'退出程序。")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

func processFile(filePath string) {
	fileNameWithExt := filepath.Base(filePath)
	fileNameWithoutExt := strings.TrimSuffix(fileNameWithExt, ".txt")

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("打开文件时出错: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	conversations := []Conversation{}

	var currentMessages []Message
	for scanner.Scan() {
		if len(currentMessages) == 2 {
			conversations = append(conversations, Conversation{Messages: append([]Message{{Role: "system", Content: fileNameWithoutExt}}, currentMessages...)})
			currentMessages = []Message{}
		}
		line := scanner.Text()
		role := "user"
		if len(currentMessages) == 1 {
			role = "assistant"
		}
		currentMessages = append(currentMessages, Message{Role: role, Content: line})
	}
	if len(currentMessages) > 0 {
		conversations = append(conversations, Conversation{Messages: append([]Message{{Role: "system", Content: fileNameWithoutExt}}, currentMessages...)})
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取文件时出错: %v\n", err)
		return
	}

	outputFilePath := strings.TrimSuffix(filePath, ".txt") + ".jsonl"
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("创建输出文件时出错: %v\n", err)
		return
	}
	defer outputFile.Close()

	for _, conv := range conversations {
		jsonData, err := json.Marshal(conv)
		if err != nil {
			fmt.Printf("序列化JSON时出错: %v\n", err)
			continue
		}
		outputFile.WriteString(string(jsonData) + "\n")
	}

	fmt.Println("转换完成，输出文件为:", outputFilePath)
}

func main() {
	for {
		filePath := readFilePath()
		if filePath == "exit" {
			fmt.Println("程序退出。")
			break
		}
		processFile(filePath)
	}
}
