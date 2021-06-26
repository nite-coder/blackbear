# 配置管理

## 介紹
用來管理系統配置，目前支援 `環境變數`, `文件檔案` 等模組， 執行順序依照加入模組的先後

* 環境變數: 讀取環境變數
* 文件檔案: 讀取檔案內容，支援 `YAML`, `JSON` 等格式

## 使用方式

範例配置文件 (app.yml)
```yml
app_id: blackbear
web:
  port: 10080
  ping: true
```

1. 設定 `文件檔案` 模組
```Go
    fileProvder := file.New()
    fileProvder.SetConfigName("config.yml") // 手動修改為讀取 config.yml, 預設讀取 app.yml 檔名
    fileProvder.SetConfigType("yaml") // 如果檔案沒有附檔名需要設定，支援 yaml
    err := fileProvder.Load() // 如果目錄下都找不到配置檔，ErrFileNotFound 會被回傳
    if err != nil {
        return err
    }

    config.AddProvider(fileProvider)
```

1. 讀取配置內容
```Go
    fileProvder := file.New()
    err := fileProvder.Load()
    if err != nil {
        return err
    }
    config.AddProvider(fileProvider)
    appID, err := config.String("app.id") // case casesentive
    fmt.Print(appID) // print: blackbear

    port, err := config.Int32("web.port", 10080) // 設定預設值，如果 "web.port" 這個 key 找不到, 就會回傳 "10080"
    fmt.Print(port) // print: 10080
```



## 更新檢查
目前是採用緩存機制，如果已經有內容被讀入就會被緩存已提升後續的效能

## RoadMap
1. File Watch
1. 配置檔的繼承模式
1. integrate remote config system

