# Simple face recognizer
[Available at](https://damp-lake-99927.herokuapp.com/)
[Available at ](https://thawing-sea-83431.herokuapp.com/)

## Photofinisih
Сервис для поиска фотографий спортсмена по номеру участника

## Setting Credentials
Setting your credentials for use by the AWS SDK can be done in a number of ways,  
but here are the recommended approaches:

Set credentials in the AWS credentials profile file on your local system, located at:

`~/.aws/credentials` on Linux, macOS, or Unix

`C:\Users\USERNAME\.aws\credentials` on Windows

This file should contain lines in the following format:

Install [go-migrate](https://github.com/golang-migrate/migrate)

`+
[default]`  
`aws_access_key_id = your_access_key_id` 
`aws_secret_access_key = your_secret_access_key`
## Clone repo
`git clone https://github.com/bearname/photofinish-frontend frontend`
## First building run command
`make download`
## Build
`make build`
## Run
`make up`
## Migration
`make migrateup`

Open http://localhost:3000

### Use case
![Use case](./docs/images/use-case.png)

### Components
![Components](./docs/images/components.png)

### Sequence Diagram
![Sequence Diagram](./docs/images/sequence-diagram.png)


### Needed 
![Sequence Diagram](./docs/Architecture%20-%20Copy.png)


Шаги при добавлении новых изображений:

1) Получить список изображений из папки dropbox
2) Сохранить список изображений в бд
Loop:
  1) Сохранить Оригинальное изображение в S3 
  2) Сжать изображение для preview, сохранить в S3 
  3) Разспознать текст из изображения
  4) Сжать изображение для mobile и сохранить в S3


```
// type Picture struct {
//    ID
//    Path
//
//    IsOriginalSaved
//    IsPreviewScaled
//    IsTextRekongized
//    IsMobileScaled
//}
```

Добавить запуск rawImageHandler (сервис по обработке не удачно обработанных изображений)
Добавить запуск знатия заблоченных ворекеров
Добавить возможность отслеживания процента выполненных задач

