# radio-t-data

Тут находятся данные про все выпуски подкаста Радио-Т. Пример файла:

```
$ cat data/episodes/689.json
{
    "number": 689,
    "file": {
        "size_bytes": 84848415,
        "length_seconds": 7067,
        "url": "http://cdn.radio-t.com/rt_podcast689.mp3",
        "md5": "ce7bc33158d630547a6ae0ed541efcfc"
    }
}
```

В папке `code` находится код, который создает эти файлы (код запускается по крону через GitHub Actions в этом репозитории)

Описание зачем это было сделано и некоторые подробности о том как это было сделано — в блоге https://ivan.bessarabov.ru/blog/radio-t-podcat-length
