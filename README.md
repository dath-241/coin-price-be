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
# Danh sách meeting minutes [`/reports`](https://github.com/dath-241/coin-price-be-go/tree/main/reports)

# Một số endpoint chính nhánh [develop](https://github.com/dath-241/coin-price-be-go/tree/develop)

### Authentication
- Đăng ký tài khoản
-- POST: ```/auth/register```
- Đăng nhập
-- POST: ```/auth/login```
- Đăng nhập bằng google
-- POST: ```/auth/google-login```
- Quên mật khẩu
-- POST: ```/auth/forgot-password```
- Đổi mật khẩu
-- POST: ```/auth/reset-password```
- Đăng xuất
-- POST: ```/auth/logout```
- Làm mới token
-- POST: ```/auth/refresh-token```

### Admin
- Lấy thông tin tất cả người dùng:
-- GET: ```/api/v1/admin/users```
- Lấy thông tin người dùng
-- GET: ```/api/v1/admin/user/:user_id```
- Xóa người dùng
-- DELETE: ```/api/v1/admin/user/:user_id```
- Cấm người dùng
-- PUT: ```/api/v1/admin/user/:_user_id/ban```
- Bỏ ban người dùng
-- PUT: ```/api/v1/admin/user/:_user_id/active```
- Lấy lịch sử thanh toán của tất cả người dùng
-- GET: ```/api/v1/admin/payment-history```
- Lấy lịch sử thanh toán của người dùng
-- GET: ```/api/v1/admin/payment-history/:user_id```

### User 
- Lấy thông tin người dùng
-- GET: ```/api/v1/user/me```
- Cập nhật thông tin người dùng
-- POST: ```/api/v1/user/me```
- Xóa tài khoản người dùng
-- DELETE: ```/api/v1/user/me```
- Đổi mật khẩu người dùng
-- PUT: ```/api/v1/user/me/change-password```
- Đổi email người dùng
-- PUT: ```/api/v1/user/me/change-email```
- Lấy lịch sử thanh toán
-- GET: ```/api/v1/user/me/payment-history```

### Coin price
- Lấy giá spot realtime
-- SUBCRIBE Websocket: ```/api/v1/spot-price/websocket```
- Lấy giá future realtime
-- SUBCRIBE Websocket: ```/api/v1/future-price/websocket```
- Lấy funding rate realtime
-- SUBCRIBE Websocket: ```/api/v1/funding-rate/websocket```
- Lấy funding rate
-- GET: ```/api/v1/funding-rate```
- Lấy kline realtime
-- SUBCRIBE Websocket: ```/api/v1/vip1/kline/websocket```
- Lấy kline
-- GET: ```/api/v1/vip1/kline```
- Lấy Market Cap
-- GET: ```/api/v1/market-stats```

### Trigger
- Tạo trigger
-- POST: ```/api/v1/vip2/alerts```
- Lấy toàn bộ trigger
-- GET: ```/api/v1/vip2/alerts```
- Lấy trigger
-- GET: ```/api/v1/vip2/alerts/:id```
- Xóa trigger
-- DELETE: ```/api/v1/vip2/alerts/:id```
- Lấy toàn bộ new/delisted symbol 
-- GET: ```/api/v1/vip2/symbol-alerts``` 
- Tạo trigger cho new/delisted symbol
-- POST: ```/api/v1/vip2/alerts/symbol```
- Bắt đầu trigger
-- POST: ```/api/v1/vip2/start-alert-checker```
- Dừng trigger
-- POST: ```/api/v1/vip2/stop-alert-checker```
- Tạo trigger cho advanced indicators
-- POST: ```/api/v1/vip3/indicators```
- Lấy toàn bộ trigger của người dùng
-- GET: ```/api/v1/users/:user_id/alerts```
- Gửi thông báo của trigger đến người dùng
-- POST: ```/api/v1/users/:user_id/alerts/notify```

### Payment
- Thanh toán nâng cấp VIP qua Momo
-- POST: ```/api/v1/payment/vip-upgrade```
- IPN (Instant Payment Notification) cho Momo
-- POST: ```/api/v1/payment/momo-callback```
- Kiểm tra trạng thái giao dịch qua Momo
-- POST: ```/api/v1/payment/status```
- Xác nhận thanh toán từ Momo và nâng cấp VIP
-- POST: ```/api/v1/payment/confirm```

#### Tham khảo thêm qua nhánh [develop](https://github.com/dath-241/coin-price-be-go/tree/develop) và các [issue](https://github.com/dath-241/coin-price-be-go/issues)
