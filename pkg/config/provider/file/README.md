# 文件檔案模組


## 讀取屬性的優先順序
有效順序依照下面排序，數字越小有效等級越高，如果檔案已經被找到將不會繼續往下找尋
1. A `/config` subdir of the working directory.
1. The working directory
1. 當前執行檔所位址下面的 `config` 目錄
1. 當前執行檔所在目錄



