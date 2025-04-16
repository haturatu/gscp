# gscp
scpのラッパーとして機能し、scpの転送ジョブを`goroutine`の並列処理によって管理します。  
  
```bash
git clone https://github.com/haturatu/gscp.git
cd gscp/cmd && go build -o gscp main.go
./gscp
```

## example
```
./gscp -d "/mnt/dest/dir" -u "taro" -h "10.1.0.1" -r "/media/src/dir" -P 110
```
