package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	//读取文件，产生流
	c := readFile()
	//单线程消费流
	output := reverseWorker(c)
	// 消费输出结果
	func(out <-chan []byte) {
		file, err := os.OpenFile("B.txt", os.O_WRONLY|os.O_TRUNC, 0600)
		defer file.Close()
		for word := range out {
			log.Println("final:", (word))
			if err != nil {
				fmt.Println(err.Error())
			} else {
				_, err = file.Write(word)
				//checkErr(err)
			}
		}
	}(output)
}

func readFile() <-chan []byte {
	output := make(chan []byte)
	go func() {
		file, err := os.Open("./A.txt")
		if err != nil {
			log.Fatal("读取文件失败！")
		}
		defer file.Close()
		buf := make([]byte, 1) //一次读取多少个字节
		bfRd := bufio.NewReader(file)
		word := make([]byte, 0) //TODO what if a work longer than 100
		for {
			n, err := bfRd.Read(buf)
			if n == 0 { //可能EOF了
				if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
					if err == io.EOF { //文件读取完毕
						output <- word
						word = make([]byte, 0)
						break; //正确退出循环，关闭channel
					}
					//TODO deal with error
				}
			}
			word = append(word, buf[0])
			if buf[0] == 32 { //遇到空格，则直接打印
				output <- word
				word = make([]byte, 0)
				continue
			}
		}
		close(output)
	}()
	return output
}

func reverseWorker(input <-chan []byte) <-chan []byte {
	output := make(chan []byte)
	go func() {
		for bytes := range input {
			output <- reverseBytes(bytes)
		}
		close(output)
	}()
	return output
}

func reverseBytes(bytes []byte) []byte {
	reverseBytes := make([]byte, len(bytes))
	copy(reverseBytes, bytes)
	for i := len(reverseBytes)/2 - 1; i >= 0; i-- {
		opp := len(reverseBytes) - 1 - i
		reverseBytes[i], reverseBytes[opp] = reverseBytes[opp], reverseBytes[i]
	}
	return reverseBytes
}
