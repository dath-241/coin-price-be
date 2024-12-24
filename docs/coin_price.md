# API Price Documentation

- [Tổng quan](#tổng-quan)
  - [Hiện thực API về giá](#i-hiện-thực-api-về-giá)
  - [Hiện thực API về Funding Rate](#ii-hiện-thực-api-về-funding-rate)
  - [Hiện thực API về Kline](#iii-hiện-thực-api-về-kline)
  - [Hiện thực API về Market Stats](#iv-hiện-thực-api-về-market-stats)
- [Clean code](#clean-code)
- [Đánh giá dự án](#đánh-giá-dự-án)

## Tổng quan
Tổng quan về xây dựng các API cho việc lấy giá (giá spot và giá future), funding rate, kline và market stats.

---
### I. Hiện thực API về giá
- **Thông tin chi tiết về Querry Params và Response:** [TRUY CẬP TẠI ĐÂY](https://drive.google.com/file/d/1LmEpkU_9MdkfQjx2VnRXrpGs_ldM53we/view?usp=drive_link)
#### 1. Lấy giá spot
- **Mô tả**: Cho phép người dùng lấy giá spot
- **Method**: GET
- **Endpoint**: `/api/v1/spot-price`
- **Query Params**: `symbol`

- **Responses**:
  - **200**: Lấy giá spot thành công
  - **400**: Query params không hợp lệ (thiếu query params hoặc thiếu symbol)
  - **500**: Lỗi server (Query params không tồn tại)

#### 2. Lấy giá future
- **Mô tả**: Cho phép người dùng lấy giá future
- **Method**: GET
- **Endpoint**: `/api/v1/future-price`
- **Query Params**: `symbol`

- **Responses**:
  - **200**: Lấy giá future thành công
  - **400**: Query params không hợp lệ (thiếu query params hoặc thiếu symbol)
  - **500**: Lỗi server (Query params không tồn tại)

#### 3. Lấy giá spot bằng websocket
- **Mô tả**: Cho phép người dùng lấy giá spot bằng cách sử dụng websocket
- **Endpoint**: `/api/v1/spot-price/websocket`
- **Query Params**: `symbol`

- **Responses**:
  - **1000**: Lấy giá spot thành công
  - **1002**: Query params không hợp lệ (thiếu query params hoặc thiếu symbol)
  - **1002**: Lỗi server (Query params không tồn tại)


#### 4. Lấy giá future bằng websocket
- **Mô tả**: Cho phép người dùng lấy giá future bằng cách sử dụng websocket
- **Endpoint**: `/api/v1/future-price/websocket`
- **Query Params**: `symbol`

- **Responses**:
  - **1000**: Lấy giá future thành công
  - **1002**: Query params không hợp lệ (thiếu query params hoặc thiếu symbol)
  - **1002**: Lỗi server (Query params không tồn tại)

---

### II. Hiện thực API về Funding Rate
- **Thông tin chi tiết về Querry Params và Response:** [TRUY CẬP TẠI ĐÂY](https://drive.google.com/file/d/1LmEpkU_9MdkfQjx2VnRXrpGs_ldM53we/view?usp=drive_link)
#### 1. Lấy funding rate
- **Mô tả**: Cho phép người dùng lấy funding rate
- **Method**: GET
- **Endpoint**: `/api/v1/funding-rate`
- **Query Params**: `symbol`

- **Responses**:
  - **200**: Lấy giá funding rate thành công
  - **400**: Query params không hợp lệ (thiếu query params hoặc thiếu symbol)
  - **400**: Lỗi server (Query params không tồn tại)

#### 2. Lấy funding rate bằng websocket
- **Mô tả**: Cho phép người dùng lấy funding rate bằng websocket
- **Endpoint**: `/api/v1/funding-rate/websocket`
- **Query Params**: `symbol`

- **Responses**:
  - **1000**: Lấy giá funding rate thành công
  - **1002**: Query params không hợp lệ (thiếu query params hoặc thiếu symbol)
  - **1002**: Lỗi server (Query params không tồn tại)

---

### III. Hiện thực API về Kline
- **Thông tin chi tiết về Querry Params và Response:** [TRUY CẬP TẠI ĐÂY](https://drive.google.com/file/d/1LmEpkU_9MdkfQjx2VnRXrpGs_ldM53we/view?usp=drive_link)
#### 1. Lấy Kline
- **Mô tả**: Cho phép người dùng lấy Kline
- **Method**: GET
- **Endpoint**: `/api/v1/vip1/kline`
- **Query Params**: `symbol`,  `interval`
- Headers: `Authorization`

- **Responses**:
  - **200**: Lấy giá Kline thành công
  - **400**: Query params không hợp lệ (thiếu query params hoặc thiếu symbol, interval)
  - **401**: Người dùng không có quyền truy cập
  - **404**: Query params không tồn tại
  - **500**: Lỗi symbol, interval sai format

#### 2. Lấy Kline bằng websocket
- **Mô tả**: Cho phép người dùng lấy Kline bằng websocket
- **Endpoint**: `/api/v1/vip1/kline/websocket`
- **Query Params**: `symbol`,  `interval`

- **Responses**:
  - **1000**: Lấy giá Kline thành công
  - **1002**: Query params không hợp lệ (thiếu query params hoặc thiếu symbol, interval)
  - **1002**: Lỗi server (Query params không tồn tại)

---

### IV. Hiện thực API về Market Stats
- **Thông tin chi tiết về Querry Params và Response:** [TRUY CẬP TẠI ĐÂY](https://drive.google.com/file/d/1LmEpkU_9MdkfQjx2VnRXrpGs_ldM53we/view?usp=drive_link)
#### 1. Lấy Market Stats bằng websocket
- **Mô tả**: Cho phép người dùng lấy Market Stats bằng websocket
- **Endpoint**: `/api/v1/market-stats`
- **Query Params**: `symbol`

- **Responses**:
  - **1000**: Lấy giá spot thành công
  - **1002**: Query params không hợp lệ (thiếu query params hoặc thiếu symbol)
  - **1002**: Lỗi server (Query params không tồn tại)

---
## Clean Code
### I. Cấu trúc thư mục
```plaintext
.
├── /price_service
│   └── /models	 # Định nghĩa các cấu trúc dữ liệu (struct)
│   └── /repository	# # Định nghĩa các interface để tương tác với database 
│   └── /routes	# Định nghĩa các endpoint của API và ánh xạ chúng tới các controller tương ứng.
│   └── /services	# Định nghĩa các logic xử lý yêu cầu (request) mà server nhận được từ client.
│   └── /utils	# Các hàm tiện ích hỗ trợ 


```
- Cấu trúc này giúp dự án dễ dàng mở rộng, phân tách rõ ràng các chức năng và dễ bảo trì.

### II. Tên biến, hàm, struct
- Đươc đặt tên theo nguyên tắc self-explanatory(tự giải thích):
- Ví dụ:
  ```go
  type ResponseSpotPrice struct {
	Symbol    string `json:"symbol"`
	Price     string `json:"price"`
	EventTime string `json:"eventTime"`
}
    ```

### III. Xử lý lỗi
- Sử dụng cấu trúc điều kiện if-else để kiểm tra và phản hồi lỗi với thông tin cụ thể. Tránh để các lỗi chưa được xử lý (unhandled errors).

---

## Đánh giá dự án

### I. Những cái làm được:
#### Lấy giá Funding Rate (`/api/v1/funding-rate` - GET)
- **Các điều kiện thành công**:
  - Thành công với các điều kiện khi symbol tồn tại

#### Lấy giá kline (`/api/v1/vip1/kline` - GET)
- Gửi email thành công khi điều kiện cảnh báo đạt.

#### Lấy giá spot (`/api/v1/spot-price` - GET)
- Thành công với các điều kiện khi symbol tồn tại

#### Lấy giá future (`/api/v1/future-price` - GET)
- Thành công với các điều kiện khi symbol tồn tại

#### Lấy websocket cho giá spot (`/api/v1/spot-price/websocket`)
- Thành công khi symbol tồn tại
- Sau mỗi 1 giây, hệ thống trả về giá kline cho người dùng

#### Lấy websocket cho giá future (`/api/v1/future-price/websocket`)
- Thành công khi symbol tồn tại
- Sau mỗi 1 giây, hệ thống trả về giá kline cho người dùng

#### Lấy websocket cho giá kline (`/api/v1/vip1/kline/websocket`)
- Thành công khi symbol tồn tại
- Xác thực token và quyền truy cập (VIP2, VIP3)
- Sau mỗi 1 giây, hệ thống trả về giá kline cho người dùng

#### Lấy websocket cho coinmarketcap (`/api/v1/market-stats`)
- Thành công khi symbol tồn tại
- Sau mỗi 15 phút, hệ thống trả về giá trị market_cap và 24h_volume cho người dùng

#### Đánh giá mức độ hoàn thành:
- RESTful API và Websocket đều được test kỹ lưỡng (Manual test với Postman và Unit Test sử dụng go test đạt trên 80% coverage).
- Handle và trả về các lỗi tương ứng.

### II. Những cái chưa làm được được:
- Các symbol chưa đồng nhất (symbol giữa coinmarket với symbol của các API còn lại).
- Chưa xử lý được trường hơp symbol người dùng nhập tồn tại nhưng api bên thứ 3 không xử lý được.
- Chưa test được các trường hợp liên quan đến lỗi của các thư viện bên thứ 3.
