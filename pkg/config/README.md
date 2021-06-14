# 配置管理

## 介紹
用來管理配置的模組，目前支援 `yaml` 的格式

## 讀取屬性的優先順序
有效順序依照下面排序，數字越小有效等級越高，如果檔案已經被找到將不會繼續往下找尋
1. OS environment variables. (all upper case)
1. A `/config` subdir of the working directory.
1. The working directory
1. 當前執行檔所位址下面的 `config` 目錄
1. 當前執行檔所在目錄



## 使用方式

範例配置文件 (app.yml)
```yml
app_id: blackbear
web:
  port: 10080
  ping: true
```

1. 讀取配置實例
```Go
    config.SetFileName("config.yml") // 預設讀取 app.yml 檔名, 手動修改為讀取 config.yml
    config.SetConfigType("yaml") // 如果檔案沒有附檔名需要設定，支援 yaml
    err := config.Load() // 如果目錄下都找不到配置檔，ErrFileNotFound 會被回傳
    if err != nil {
        return err
    }

    cfg := config.Cfg() // singleton
```

1. 讀取配置內容
```Go
    err := config.Load()
    if err != nil {
        return err
    }
    appID, err := config.String("app.id")
    fmt.Print(appID) // print: blackbear

    port, err := config.String("web.port", "10080") // 設定預設值，如果 "web.port" 這個 key 找不到, 就會回傳 "10080"
    fmt.Print(port) // print: 10080
```

## 更新檢查
目前是採用緩存機制，如果已經有內容被讀入就會被緩存已提升後續的效能

## RoadMap
1. 配置檔的繼承模式
1. Apollo 等工具整合
1. File Watch
