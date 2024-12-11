# RESTful API Documentation

## I. RESTful API

### 1. Lấy giá Funding Rate
**Method:** `GET`  
**Endpoint:** `/api/v1/funding-rate`  
**Params:** `symbol`

**Cơ chế:**  
Fetch API từ: `https://fapi.binance.com/fapi/v1/fundingInfo`. Khi người dùng gọi API đến server, server sẽ:
- Lấy thông tin `symbol` từ người dùng nhập vào.
- Truy xuất thông tin qua API từ bên thứ 3 (Binance).

**Xử lý:**
- **Trường hợp thiếu hoặc không hợp lệ:**
  - Thiếu trường `symbol`: Trả về **400 Bad Request** với thông báo lỗi phù hợp.
  - `symbol` không tồn tại: Trả về **400 Bad Request** với thông báo lỗi phù hợp.
- **Trường hợp hợp lệ:** Trả về **200 OK** và các thông tin liên quan đến `symbol`, bao gồm:
  - `fundingRate`: Giá trị funding rate của `symbol` nhập vào.
  - `fundingCountDown`: Thời gian đếm ngược để reset.
  - `eventTime`: Thời gian lúc gọi API.
  - `adjustedFundingRateCap`: Mức tối đa mà tỷ lệ tài trợ điều chỉnh có thể đạt được.
  - `adjustedFundingRateFloor`: Giới hạn tối thiểu (floor) cho tỷ lệ tài trợ điều chỉnh.
  - `fundingIntervalHours`: Khoảng thời gian giữa các lần thanh toán tỷ lệ tài trợ (thường là 8 giờ).
- **Trường hợp 1**: Thiếu symbol
![Screenshot 2024-12-07 201040](https://github.com/user-attachments/assets/bdd2f08e-6260-48af-8537-de1643c2fd72)
- **Trường hợp 2**: symbol không tồn tại
![Screenshot 2024-12-07 201507](https://github.com/user-attachments/assets/a77f7421-410a-47ea-a0e6-cf7bb53bf80c)
- **Trường hợp 3**: symbol tồn tại
![Screenshot 2024-12-07 201545](https://github.com/user-attachments/assets/e3014bc9-271c-4fd4-b9c5-eaa542492900)


---

### 2. Lấy giá Kline
**Method:** `GET`  
**Endpoint:** `/api/vip1/kline`  
**Headers:** `Authorization`  
**Params:** `symbol`, `interval`

**Cơ chế:**  
Fetch API từ: `https://fapi.binance.com/fapi/v1/klines`. Khi người dùng gọi API đến server, server sẽ:
- Lấy thông tin `symbol`, `interval` từ người dùng nhập vào.
- Truy xuất thông tin qua API từ bên thứ 3 (Binance).

**Quy định truy cập:** API này yêu cầu tài khoản có **role từ VIP2 trở lên**. Các tài khoản có role thấp hơn sẽ không được truy cập.

**Xử lý:**
- **Trường hợp thiếu hoặc không hợp lệ:**
  - Thiếu `symbol`: Trả về **400 Bad Request**.
  - Thiếu `interval`: Trả về **400 Bad Request**.
  - `symbol` không tồn tại: Trả về **404 Not Found**.
  - `interval` sai format: Trả về **400 Bad Request**.
- **Trường hợp hợp lệ:** Trả về **200 OK** và thông tin Kline.
- **Trường hợp 1**: Người dùng không có quyền truy cập
![Screenshot 2024-12-07 203736](https://github.com/user-attachments/assets/32a445e3-825e-4e83-8f43-c50f8513b736)
- **Trường hợp 2**: Người dùng có quyền truy cập, nhưng thiếu symbol
![Screenshot 2024-12-07 203826](https://github.com/user-attachments/assets/16f830a2-dd60-4f3b-b6e5-0eceef48a67d)
- **Trường hợp 3**: Người dùng có quyền truy cập, nhưng thiếu interval
![Screenshot 2024-12-07 203906](https://github.com/user-attachments/assets/29caf29a-de4c-4499-962b-b7cfeed52aee)
- **Trường hợp 4**: Người dùng có quyền truy cập, symbol tồn tại, interval đúng format
![Screenshot 2024-12-07 204035](https://github.com/user-attachments/assets/4c9374d7-d04c-4b50-b441-503494f01e43)
- **Trường hợp 5**: Người dùng có quyền truy cập, symbol không tồn tại
![Screenshot 2024-12-07 204200](https://github.com/user-attachments/assets/25d6a9f8-88ae-42c2-ac91-867008e2d47e)
- **Trường hợp 6**: Người dùng có quyền truy cập, symbol tồn tại, interval sai format
![Screenshot 2024-12-07 204255](https://github.com/user-attachments/assets/b92ff1e4-81f1-4373-a93d-c5512c6c8982)

---

## WebSocket API

### 1. Lấy giá Funding Rate
**Endpoint:** `/funding-rate/websocket`  
**Query Params:** `symbol`

**Cơ chế:**  
Fetch từ: `wss://stream.binance.com/stream?streams=%s@markPrice@1s`

**Xử lý:**
- **TH1: Thiếu `symbol`**  
  Sau 5 giây, socket sẽ tự động đóng và trả về:
  - Mã lỗi: `1002`
  - Message: `Symbol error`
  ![Screenshot 2024-12-07 213129](https://github.com/user-attachments/assets/6d765b35-0613-45b0-a246-97f39383b733)


- **TH2: `symbol` không tồn tại**  
  Nếu không nhận được phản hồi từ Binance sau 5 giây, hệ thống sẽ:
  - Trả về mã lỗi: `1002`
  - Message: `Symbol error`
  - Đóng socket.
![Screenshot 2024-12-07 213910](https://github.com/user-attachments/assets/93d6a7aa-db6f-4fc4-8e07-6080ed67835f)

- **TH3: `symbol` tồn tại**  
  Nếu Binance trả về thông tin hợp lệ, server sẽ xử lý và gửi lại cho người dùng sau mỗi 1 giây.
![Screenshot 2024-12-07 214108](https://github.com/user-attachments/assets/b6c13682-8c2e-4394-a02f-b9a971ec3d90)
  Thông tin trả về sẽ có định dạng như sau:
  ```json
  {
    "eventTime": "2021-12-04 15:53:23",
    "fundingCountDown": "00:06:37",
    "fundingRate": "0.00039888",
    "symbol": "QTUMUSDT"
  }
  ```
  
---

### 2. Lấy giá Kline
**Endpoint:** `/kline/websocket`  
**Headers:** `Authorization`  
**Query Params:** `symbol`

**Cơ chế:**  
Fetch từ: `wss://stream.binance.com/stream?streams=%s@kline_1s`

**Quy định truy cập:** Yêu cầu tài khoản từ **VIP2 trở lên**.

**Xử lý:**
- **TH1: Không có quyền truy cập**  
  Nếu role thấp hơn VIP2, socket sẽ tự động đóng sau 5 giây. (Hiện tại đang bị lỗi)
  
- **TH2: Thiếu `symbol`**  
  Socket sẽ tự động đóng sau 5 giây, trả về mã lỗi 1002.
![Screenshot 2024-12-07 225511](https://github.com/user-attachments/assets/ea2afc85-c89d-465a-80a8-f6228df8189a)

- **TH3: `symbol` không tồn tại**  
  Socket sẽ tự động đóng sau 5 giây, trả về mã lỗi.
![Screenshot 2024-12-07 225623](https://github.com/user-attachments/assets/0c578de1-ddf3-4982-a1b9-a907d1e366d4)

- **TH4: `symbol` tồn tại**  
  Server trả về response sau 1 giây.
![Screenshot 2024-12-07 225706](https://github.com/user-attachments/assets/19ed2f6c-f7bb-42a8-b4c8-b0ed41d93901)
  Thông tin trả về sẽ có định dạng như sau:
  ```json
  {
    "baseAssetVolume": "0.06551000",
    "closeTime": "2024-12-04 16:28:22",
    "eventTime": "2024-12-04 16:28:23",
    "highPrice": "94984.14000000",
    "lowPrice": "94984.13000000",
    "openPrice": "94984.14000000",
    "quoteAssetVolume": "6222.4105860",
    "startTime": "2024-12-04 16:28:22",
    "symbol": "BTCUSDT",
    "takerBuyBaseVolume": "0.02293000",
    "takerBuyQuoteVolume": "2177.98633020"
  }
  ```

---

### 3. CoinMarketCap
**Endpoint:** `/market-stats`  
**Query Params:** `symbol`

**Cơ chế:**  
Fetch từ: `https://api.coingecko.com/api/v3/coins/%s`

**Xử lý:**
- **TH1: Thiếu `symbol`**  
  Server sẽ đóng socket và trả về:
  - Mã lỗi: `1000`
  - Message: `Symbol missing or invalid`
![Screenshot 2024-12-08 000107](https://github.com/user-attachments/assets/25c374d6-c98e-4cce-a077-d493c5aa48ce)

- **TH2: `symbol` không tồn tại**  
  Server sẽ đóng socket và trả về:
  - Mã lỗi: `1000`
  - Message: `Symbol missing or invalid`
![Screenshot 2024-12-08 000239](https://github.com/user-attachments/assets/43fcfb95-d898-46fc-8eca-fddb0b54dadb)


- **TH3: `symbol` tồn tại**  
  Server trả về response cho người dùng và tự động cập nhật sau 15 phút.
  ![image](https://github.com/user-attachments/assets/d1b65ca5-ae91-484e-980d-f3fbe4744fd5)
  
  ```json
  {
    "symbol": "btc",
    "market_cap": 1973049246845,
    "24h_volume": 84464126237
  }
  ```
