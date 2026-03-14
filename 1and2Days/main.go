package main
import (
	"fmt"
	"errors"
)

func greet(name string, age int) (string, error) {
	if age < 0 {
		return "", errors.New("Возраст не может быть отрицательным")
	}
	text := fmt.Sprintf("Привет %s, тебе %d лет", name, age)

	return text, nil
}

var population = map[string]int{
	"Москва": 12506468,
	"Санкт-Петербург": 5351935,
	"Новосибирск": 1620162,
	"Екатеринбург": 1493749,
	"Нижний Новгород": 1250615,
}

func getPopu(city string) (string, int, error){
	pop, ok := population[city]
	if !ok {
		return "", 0, errors.New("города нету в бд")
	}
	return city, pop, nil
}


func main(){
	msg, err := greet("Алексей", 30)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println(msg)

	city, pop, err := getPopu("Москва")
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println("Население города: ", pop, "Город: ", city)
}


