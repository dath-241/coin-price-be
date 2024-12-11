## Chức Năng Hoàn Thành

### Tạo Cảnh Báo Giá (`/api/v1/vip2/alerts` - POST)
- **Các điều kiện thành công**:
  - Spot/Future giá đạt mốc cụ thể (>=, <=).
  - Giá nằm trong khoảng (min_range, max_range).
  - Funding rate vượt mốc hoặc thay đổi interval.
- **Xử lý request**:
  - Kiểm tra đầy đủ các trường dữ liệu (symbol, threshold, condition...).
  - Xác thực token và quyền truy cập (VIP2/VIP3).
  - Giới hạn tối đa 5 cảnh báo mỗi người dùng.
- **Tích hợp điều kiện snooze**:
  - "Only once", "Once per duration", "At specific time".

### Gửi Thông Báo Đến Người Dùng (`/api/v1/users/:id/alerts/notify` - POST)
- Gửi email thành công khi điều kiện cảnh báo đạt.

### Kiểm Tra Cảnh Báo (`/api/v1/vip2/start-alert-checker` - POST)
- Hệ thống tự động kiểm tra các điều kiện cảnh báo.
- Kích hoạt thông báo nếu điều kiện hợp lệ.

### Dừng Kiểm Tra Cảnh Báo (`/api/v1/vip2/stop-alert-checker` - POST)
- Ngừng kiểm tra điều kiện và gửi thông báo.

### Lấy Danh Sách Cảnh Báo (`/api/v1/vip2/alerts` - GET)
- Trả về danh sách cảnh báo của người dùng hiện tại.
- Xử lý lỗi: thiếu Authorization, token không hợp lệ.

### Lấy Thông Tin Cảnh Báo Cụ Thể (`/api/v1/vip2/alerts/:id` - GET)
- Lấy thông tin một cảnh báo dựa trên ID.
- Xử lý lỗi: ID không hợp lệ hoặc không tìm thấy cảnh báo.

### Xóa Cảnh Báo (`/api/v1/vip2/alerts/:id` - DELETE)
- Xóa cảnh báo dựa trên ID.
- Xử lý lỗi tương tự như GET một cảnh báo.

---

## Chức Năng Chưa Hoàn Thành

- Cảnh báo dựa trên **indicator** như MA, EMA, BollingerBands.
- Kiểm tra và gửi thông báo khi **symbol newlisting** hoặc **delisting**.

---

## Kịch Bản Demo - Mô Tả Luồng Chạy

### Kịch Bản Chính
Tập trung vào các chức năng hệ thống cảnh báo giá, bao gồm:
1. Tạo cảnh báo với các điều kiện khác nhau.
2. Kích hoạt kiểm tra cảnh báo.
3. Gửi thông báo khi cảnh báo thỏa điều kiện.
4. Xem danh sách và chi tiết cảnh báo.
5. Xóa cảnh báo khi không cần thiết.

### Chi Tiết Các Bước
#### **Bước 1: Tạo Cảnh Báo**
- Gửi yêu cầu `POST` đến endpoint `/api/v1/vip2/alerts` với các payload khác nhau:
  - Cảnh báo giá spot >= 60,000.
  - Cảnh báo giá spot trong khoảng 40,000 - 100,000.

#### **Bước 2: Lấy Danh Sách Cảnh Báo**
- Gửi yêu cầu `GET` đến endpoint `/api/v1/vip2/alerts`.

#### **Bước 3: Kích Hoạt Kiểm Tra Cảnh Báo**
- Gửi yêu cầu `POST` đến endpoint `/api/v1/vip2/start-alert-checker`.

#### **Bước 4: Gửi Thông Báo**
- Khi điều kiện cảnh báo đạt (ví dụ: giá BTCUSDT >= 60,000), hệ thống tự động gửi thông báo qua email.

#### **Bước 5: Xem Chi Tiết Một Cảnh Báo**
- Gửi yêu cầu `GET` đến endpoint `/api/v1/vip2/alerts/:id`.

#### **Bước 6: Xóa Một Cảnh Báo**
- Gửi yêu cầu `DELETE` đến endpoint `/api/v1/vip2/alerts/:id`.
- 
## Mô tả các API
### **TẠO CẢNH BÁO GIÁ**


**1. Description:**

- Endpoint:
`/api/v1/vip2/alerts (POST)`

- Request Body:  
```
{
  "symbol": "string",
  "price": "float",
  "condition": "string",
  "threshold": "float",
  "is_active": "boolean",
  "notification_method": "string",
  "type": "string",
  "symbols": ["string"],
  "frequency": "string",
  "snooze_condition": "string",
  "max_repeat_count": "integer",
  "next_trigger_time": "datetime (ISO 8601 format, e.g., 2023-12-08T12:00:00Z)",
  "repeat_count": "integer",
  "message": "string",
  "last_fundingrate_interval": "string",
  "min_range": "float",
  "max_range": "float"
}

```
**2. Demo:**
- Success:
```
{
    "symbol": "BTCUSDT",
    "threshold": 60000,
    "condition": ">=",
    "is_active": true,
    "notification_method": "email",
    "type": "spot",
    "snooze_condition": "Only once",
    "max_repeat_count": 5
}
```
![image](https://github.com/user-attachments/assets/fab6fb97-349b-4177-9f70-020499bc3f69)
```
{
    "symbol": "BTCUSDT",
    "min_range": 30000,
    "max_range": 60000,
    "condition": ">=",
    "is_active": true,
    "notification_method": "email",
    "type": "spot",
    "snooze_condition": "Only once",
    "max_repeat_count": 5
}
```
![image](https://github.com/user-attachments/assets/08adedd5-ea65-4260-bc9a-65bd892b9ad5)

- Invalid request body:
![image](https://github.com/user-attachments/assets/48c31d96-b798-4106-b0c6-a6736781ca8e)
- Authorization header required:
![image](https://github.com/user-attachments/assets/09e844c8-3c7c-4ff7-ae0d-821b5fa7dfb1)
- Invalid token:
![image](https://github.com/user-attachments/assets/3d17e98e-35ce-4590-8bd2-2c2b55c78406)
- Maximum alert limit reached: (max is 5)
![image](https://github.com/user-attachments/assets/d9fee8c9-660e-4c26-8fa5-0dc74b59e32e)
- Access forbidden: insufficient role: (not vip2 or vip3)
![image](https://github.com/user-attachments/assets/41f096d4-063b-47d2-907e-29985b4aa1f4)

### **TẠO THÔNG BÁO ĐẾN NGƯỜI DÙNG**
**1. Description:**

- Endpoint:
`api/v1/users/:id/alerts/notify (POST)`

**2. Demo:**
- Success:
![image](https://github.com/user-attachments/assets/92433955-69a3-4513-bf85-b3ea13bd8435)

### **BẮT ĐẦU KIỂM TRA CẢNH BÁO**
**1. Description:**
Nếu thỏa điều kiện thì gửi thông báo.

- Endpoint:
`/api/v1/vip2/start-alert-checker (POST)`

**2. Demo:**
- Success:
![image](https://github.com/user-attachments/assets/41df5879-cb17-4bb2-9c19-854d3ef31eeb)

### **DỪNG KIỂM TRA CẢNH BÁO**
**1. Description:**
Dừng kiểm tra cảnh báo và dừng gửi thông báo

- Endpoint:
`/api/v1/vip2/stop-alert-checker (POST)`

**2. Demo:**
- Success:
![image](https://github.com/user-attachments/assets/185ba313-1cd6-43ff-a99b-f68f89c86a5d)

### **LẤY TẤT CẢ CẢNH BÁO GIÁ CỦA NGƯỜI DÙNG**
**1. Description:**

- Endpoint:
`/api/v1/vip2/alerts (GET)`

**2. Demo:**
- Success:
![image](https://github.com/user-attachments/assets/8d3ce34a-2c67-4f53-b1c2-3cb0083d7547)
- Authorization header required:
![image](https://github.com/user-attachments/assets/bb0e0a08-6d96-4f56-9e9e-e7e7e29cca7d)
- Invalid token:
![image](https://github.com/user-attachments/assets/59550dad-4b0c-48e7-b41f-f228d555cf79)

### **LẤY 1 CẢNH BÁO GIÁ CỤ THỂ**

**1. Description:**

- Endpoint:
`/api/v1/vip2/alerts/:id (GET)`

**2. Demo:**
- Success:
![image](https://github.com/user-attachments/assets/1076d197-8fc3-47ee-836e-cdacb76e4e71)
- Authorization header required:
![image](https://github.com/user-attachments/assets/41888ba8-7d21-43d1-b909-7c3fcb711f48)
- Invalid token:
![image](https://github.com/user-attachments/assets/44d5703e-af52-4e43-b3c1-e39e1dc4be80)
- Invalid alert ID:
![image](https://github.com/user-attachments/assets/f09b0001-742c-4478-abbd-b24105b45f64)
- Alert not found:
![image](https://github.com/user-attachments/assets/85cfed59-e997-41dd-b2c9-b3b6d1dc08d3)


### **XÓA CẢNH BÁO GIÁ**

**1. Description:**

- Endpoint:
`/api/v1/vip2/alerts/:id (DELETE)`

**2. Demo:**
- Success:
![image](https://github.com/user-attachments/assets/00f9d741-f5ba-4b9e-b1dc-5ba69066eee6)
- Authorization header required:
![image](https://github.com/user-attachments/assets/8e755655-e779-4fb8-916b-f7d26f03acf2)
- Invalid token:
![image](https://github.com/user-attachments/assets/49527b79-ccf9-4e98-87aa-0bbe5a34a2e4)
- Invalid alert ID:
![image](https://github.com/user-attachments/assets/975a493e-7e48-4fd7-8ef8-987b0cb468b7)
- Alert not found:
![image](https://github.com/user-attachments/assets/6be229f2-40f3-466f-921e-ea9b443cc21c)
# *TRIGGER CONDITION*

### **Giá Spot/Future tăng/giảm hơn 1 mốc nào đó, vào/ra range nào đó**
- Create 4 alert:
**1. Giá spot >= 60000**
```
{
    "symbol": "BTCUSDT",
    "threshold": 60000,
    "condition": ">=",
    "is_active": true,
    "notification_method": "email",
    "type": "spot",
    "snooze_condition": "Only once",
    "max_repeat_count": 5
}
```
**2. Giá spot >= 40000 và <= 100000**
```
{
    "symbol": "BTCUSDT",
    "min_range": 40000,
    "max_range": 100000,
    "condition": ">=",
    "is_active": true,
    "notification_method": "email",
    "type": "spot",
    "snooze_condition": "Only once",
    "max_repeat_count": 5
}
```
**3. Giá future >= 60000**
```
{
    "symbol": "BTCUSDT",
    "threshold": 60000,
    "condition": ">=",
    "is_active": true,
    "notification_method": "email",
    "type": "future",
    "snooze_condition": "Only once",
    "max_repeat_count": 5
}
```
**4. Giá future >= 40000 và <= 100000**
```
{
    "symbol": "BTCUSDT",
    "min_range": 40000,
    "max_range": 100000,
    "condition": ">=",
    "is_active": true,
    "notification_method": "email",
    "type": "future",
    "snooze_condition": "Only once",
    "max_repeat_count": 5
}
```
- Notify:
![image](https://github.com/user-attachments/assets/92433955-69a3-4513-bf85-b3ea13bd8435)
![image](https://github.com/user-attachments/assets/41df5879-cb17-4bb2-9c19-854d3ef31eeb)
![image](https://github.com/user-attachments/assets/7e98e20b-fd69-464f-84c6-85cb351b85f3)
![image](https://github.com/user-attachments/assets/185ba313-1cd6-43ff-a99b-f68f89c86a5d)


### **Mức độ chênh giá giữa spot/future so với 1 mốc nào đó, vào/ra range nào đó**
- Create Alert:
```
{
    "symbol": "BTCUSDT",
    "threshold": 1,
    "condition": ">=",
    "is_active": true,
    "notification_method": "email",
    "type": "price_difference",
    "snooze_condition": "Only once",
    "max_repeat_count": 5
}
```
- Notify:
![image](https://github.com/user-attachments/assets/92433955-69a3-4513-bf85-b3ea13bd8435)
![image](https://github.com/user-attachments/assets/41df5879-cb17-4bb2-9c19-854d3ef31eeb)
![image](https://github.com/user-attachments/assets/b1da63dd-d430-42ab-b462-369bb4363565)
![image](https://github.com/user-attachments/assets/185ba313-1cd6-43ff-a99b-f68f89c86a5d)

### **Funding rate tăng/giảm hơn 1 mốc nào đó, vào/ra range nào đó**
- Create Alert:
```
{
    "symbol": "BTCUSDT",
    "threshold": 0.00001,
    "condition": ">=",
    "is_active": true,
    "notification_method": "email",
    "type": "funding_rate",
    "snooze_condition": "Only once",
    "max_repeat_count": 5
}
```
- Notify:
![image](https://github.com/user-attachments/assets/92433955-69a3-4513-bf85-b3ea13bd8435)
![image](https://github.com/user-attachments/assets/41df5879-cb17-4bb2-9c19-854d3ef31eeb)
![image](https://github.com/user-attachments/assets/5af96502-3f05-411c-bb57-7aa414894063)
![image](https://github.com/user-attachments/assets/185ba313-1cd6-43ff-a99b-f68f89c86a5d)


### **Funding rate Interval thay đổi**
- Create Alert:  
```
{
    "symbol": "BTCUSDT",
    "condition": ">=",
    "is_active": true,
    "notification_method": "email",
    "type": "funding_rate_interval",
    "last_fundingrate_interval": "9h0m0s",
    "snooze_condition": "Only once",
    "max_repeat_count": 5
}
```
- Notify:
![image](https://github.com/user-attachments/assets/92433955-69a3-4513-bf85-b3ea13bd8435)
![image](https://github.com/user-attachments/assets/41df5879-cb17-4bb2-9c19-854d3ef31eeb)
![image](https://github.com/user-attachments/assets/1246b7f2-58bd-4731-9581-b2ef4398e51d)
![image](https://github.com/user-attachments/assets/185ba313-1cd6-43ff-a99b-f68f89c86a5d)

# *SNOOZE CONDITION*

### **One-time**
- Create Alert:
```
{
    "symbol": "BTCUSDT",
    "threshold": 60000,
    "condition": ">=",
    "is_active": true,
    "notification_method": "email",
    "type": "spot",
    "snooze_condition": "Only once",
    "max_repeat_count": 5
}
```
- Notify:
![image](https://github.com/user-attachments/assets/92433955-69a3-4513-bf85-b3ea13bd8435)
![image](https://github.com/user-attachments/assets/41df5879-cb17-4bb2-9c19-854d3ef31eeb)
![image](https://github.com/user-attachments/assets/7e98e20b-fd69-464f-84c6-85cb351b85f3)
![image](https://github.com/user-attachments/assets/185ba313-1cd6-43ff-a99b-f68f89c86a5d)


### **Once in duration**
- Create Alert:
```
{
    "symbol": "BTCUSDT",
    "threshold": 60000,
    "condition": ">=",
    "is_active": true,
    "notification_method": "email",
    "type": "spot",
    "snooze_condition": "Once per 10 seconds",
    "max_repeat_count": 5
}
```
- Notify:
![image](https://github.com/user-attachments/assets/92433955-69a3-4513-bf85-b3ea13bd8435)
![image](https://github.com/user-attachments/assets/41df5879-cb17-4bb2-9c19-854d3ef31eeb)
![image](https://github.com/user-attachments/assets/2fa987e1-8778-4191-b63d-b4eccad3f873)
![image](https://github.com/user-attachments/assets/185ba313-1cd6-43ff-a99b-f68f89c86a5d)


### **At specific time**
- Create Alert:  
```
{
    "symbol": "BTCUSDT",
    "threshold": 50000,
    "condition": ">=",
    "is_active": true,
    "notification_method": "email",
    "type": "spot",
    "snooze_condition": "Once time",
    "next_trigger_time": "2024-11-15T22:29:00+07:00",
    "max_repeat_count": 5
}
```
- Notify:
![image](https://github.com/user-attachments/assets/92433955-69a3-4513-bf85-b3ea13bd8435)
![image](https://github.com/user-attachments/assets/41df5879-cb17-4bb2-9c19-854d3ef31eeb)
![image](https://github.com/user-attachments/assets/9cc19732-56cb-471e-8635-8ba355bce777)
![image](https://github.com/user-attachments/assets/185ba313-1cd6-43ff-a99b-f68f89c86a5d)





