// package main 
// import (
// 	"bytes"
// 	"fmt"
// 	"encoding/json"
// 	"io"
// 	"net/http"
// 	"os"

// 	"github.com/gin-gonic/gin"
//     "github.com/joho/godotenv"
// )

// type GemeniReq struct{
// 	Contents []Content `json:"contents"`
// }

// type Content struct {
// 	Parts []Part `json:"parts"`
// }

// type Part struct {
// 	Text string `json:"text"`
// }

// type GemeniRes struct {
// 	Candidates []struct {
// 		Content struct {
// 			Parts []Part `json:"parts"`
// 		} `json:"content"`
// 	} `json:"candidates"`
// }

// func askGemeni(prompt string) (string, error){
// 	apiKey := os.Getenv("GEMINI_API_KEY")
	
// 	url := fmt.Sprintf(
// 	    "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s",
//         apiKey,
// 	)

// 	reqBody := GemeniReq{
// 		Contents: []Content{
// 			{Parts: []Part{{Text: prompt}}},	
// 		},
// 	}
// 	jsonData, err := json.Marshal(reqBody)
// 	if err != nil {
// 		return "", err
// 	}

// 	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
// 	if err != nil{
// 		fmt.Println("Ошибка http.Post:", err) 
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", err
// 	}
// 	fmt.Println("Ответ от Gemini:", string(body))

// 	var geminiResp GemeniRes
// 	err = json.Unmarshal(body, &geminiResp)
// 	if err != nil {
// 		return "", fmt.Errorf("Пустой ответ gemini")
// 	}
// 	if len(geminiResp.Candidates) == 0 {
//     	return "", fmt.Errorf("Gemini вернул пустой ответ: %s", string(body))
// 	}
// 	return geminiResp.Candidates[0].Content.Parts[0].Text,nil
// }

// func main(){
// 	godotenv.Load()
// 	fmt.Println("KEY:", os.Getenv("GEMINI_API_KEY"))
// 	r := gin.Default()

// 	r.POST("/ask", func(c *gin.Context){
// 		var body struct {
// 			Question string `json:"question" binding:"required"`
// 		}
// 		if err := c.ShouldBindJSON(&body); err != nil {
// 			c.JSON(400, gin.H{"error": "введите ваш запрос"})
// 			return
// 		}

// 		answer, err := askGemeni(body.Question)
// 		if err != nil {
// 			c.JSON(500, gin.H{"error": err.Error()})
//             return
// 		}
		
// 		c.JSON(200, gin.H{"answer": answer})
// 	})
// 	fmt.Println("Сервер запущен на порту 8080")
//     r.Run(":8080")
// }

package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_"github.com/lib/pq"
)


type GemeniReq struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GemeniRes struct {
	Candidates []struct {
		Content struct {
			Parts []Part `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

var db *sql.DB

func initDB() {
	connStr := "postgres://postgres:910200@localhost:5432/finaltest?sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		chat_id TEXT,
		role TEXT,
		content TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

func getHistory(chatID string) ([]Content, error) {
	rows, err := db.Query(`
		SELECT content FROM messages
		WHERE chat_id=$1
		ORDER BY created_at ASC
		LIMIT 20
	`, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []Content

	for rows.Next() {
		var text string
		rows.Scan(&text)

		contents = append(contents, Content{
			Parts: []Part{{Text: text}},
		})
	}

	return contents, nil
}

func saveMessage(chatID, role, text string) {
	_, err := db.Exec(`
		INSERT INTO messages (chat_id, role, content)
		VALUES ($1, $2, $3)
	`, chatID, role, text)

	if err != nil {
		log.Println("DB error:", err)
	}
}

func askGemeni(contents []Content) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")

	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s",
		apiKey,
	)

	reqBody := GemeniReq{Contents: contents}

	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var geminiResp GemeniRes
	json.Unmarshal(body, &geminiResp)

	if len(geminiResp.Candidates) == 0 {
		return "", fmt.Errorf("empty response: %s", string(body))
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

func main() {
	godotenv.Load()
	initDB()

	r := gin.Default()

	r.POST("/ask", func(c *gin.Context) {
		var body struct {
			Question string `json:"question"`
			ChatID   string `json:"chat_id"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		history, _ := getHistory(body.ChatID)

		history = append(history, Content{
			Parts: []Part{{Text: body.Question}},
		})

		answer, err := askGemeni(history)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		saveMessage(body.ChatID, "user", body.Question)
		saveMessage(body.ChatID, "model", answer)

		c.JSON(200, gin.H{"answer": answer})
	})

	fmt.Println("Server running on :8080")
	r.Run(":8080")
}
