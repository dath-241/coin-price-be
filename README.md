# Đề tài: Hệ thống thông tin thị trường tiền điện tử
Hệ thống hỗ trợ việc theo dõi thông tin thị trường (giá tiền, phí), các giao dịch được hình thành, tra cứu thống kê các giao dịch sử dụng các API từ Binance, CoinMarketCap, CoinGecko 
# Danh sách thành viên
|    MSSV   |         Họ tên         |  Role  |
|:---------:|:----------------------:|:------:|
|  2211024  |      Lê Ngọc Hiền      | PO     |
|  2211101  |   Nguyễn Thanh Hoàng   | Dev    |
|  2210768  |    Nguyễn Văn Đoàn     | Dev    |
|  2211144  | Nguyễn Trịnh Ngọc Huân | Dev    |
|  2212962  |     Trần Quang Tác     | Dev    |
|  2210871  |     Quách Khải Hào     | Dev    |
|  2212922  |    Nguyễn Quang Sáng   | Dev    |
# Cấu trúc thư mục `/docs`

# Cấu trúc thư mục `/reports`

# Hướng dẫn sử dụng
- Clone repo về máy:
```basg
gh repo clone dath-241/coin-price-be-go
```
- Mở Docker Desktop
- Chạy docker-compose:
```basg
docker-compose up -d
```
- Server đã sẵn sàng, gọi API mẫu (POSTMAN...):
```basg
http://localhost:8080/api/v1/funding-rate?symbol=BTCUSDT
```
