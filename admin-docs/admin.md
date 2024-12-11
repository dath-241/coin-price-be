# API Admin Documentation


## Tổng quan
Tổng quan về xây dựng các API cho việc xác thực, quản lý người dùng và thanh toán.

---
### I. Hiện thực API về xác thực Authentication

#### 1. Đăng ký tài khoản người dùng
- **Mô tả**: Cho phép người dùng đăng ký bằng tên người dùng, mật khẩu, email và số điện thoại (tùy chọn).
- **Method**: POST
- **Endpoint**: `/api/v1/auth/register`
- **Request Body**:
```json
{
  "email": "user@example.com",
  "password": "hashed_password_here",
  "profile": {
    "avatar_url": "https://example.com/avatar.png",
    "bio": "Software developer based in NY",
    "date_of_birth": "1995-05-15",
    "full_name": "John Doe",
    "phone_number": "+84911123456"
  },
  "username": "johndoe"
}
```
- **Responses**:
  - **201**: Đăng ký tài khoản thành công
  - **400**: Request body không hợp lệ
  - **409**: Email, username, phone number đã tồn tại
  - **500**: Lỗi server

#### 2. Đăng nhập
- **Mô tả**: Xác thực người dùng theo tên người dùng hoặc email và trả về mã thông báo.
- **Method**: POST
- **Endpoint**: `/api/v1/auth/login`
- **Request Body**:
```json
{
  "username": "string",
  "password": "string"
}
```
- **Responses**:
  - **200**: Login thành công, trả về token
  - **400**: Request body không hợp lệ
  - **401**: Mật khẩu hoặc username không chính xác
  - **403**: Tài khoản bị cấm
  - **500**: Lỗi server

#### 3. Đăng nhập qua Google
- **Mô tả**: cho phép người dùng xác thực bằng tài khoản Google của họ. Frontend gửi 1 Google ID Token, được xác minh bởi phía backend để tạo hoặc xác thực người dùng.
- **Method**: POST
- **Endpoint**: `/api/v1/auth/google-login`
- **Parameters**: `id_token`
- **Responses**:
  - **200**: Login thành công, trả về token
  - **401**: `id_token` không hợp lệ
  - **403**: Tài khoản bị cấm
  - **500**: Lỗi server

#### 4. Quên mật khẩu
- **Mô tả**: Cho phép người dùng quên mật khẩu bằng cách gửi mã OTP đặt lại mật khẩu đến địa chỉ email của người dùng.
- **Method**: POST
- **Endpoint**: `/api/v1/auth/forgot-password`
- **Request Body**:
```json
{
  "email": "string"
}
```
- **Responses**:
  - **200**: Mã OTP được gửi đến email để đặt lại mật khẩu thành công.
  - **400**: Email không đúng định dạng
  - **404**: Không tìm thấy user với email cung cấp
  - **500**: Lỗi server

#### 5. Đặt lại mật khẩu
- **Mô tả**: cho phép người dùng đặt lại mật khẩu bằng OTP hợp lệ và mật khẩu mới. OTP được kiểm tra tính xác thực và hết hạn trước khi cập nhật mật khẩu.
- **Method**: POST
- **Endpoint**: `/api/v1/auth/reset-password`
- **Request Body**:
```json
{
  "otp": "string",
  "new_password": "string"
}
```
- **Responses**:
  - **200**: xác thực otp và đặt lại mật khẩu thành công.
  - **400**: Mật khẩu không đúng định dạng
  - **401**: OTP không hợp lệ
  - **500**: Lỗi server

#### 6. Refresh Token
- **Mô tả**: cho phép người dùng làm mới token của họ bằng token cũ hợp lệ.
- **Method**: POST
- **Endpoint**: `/api/v1/auth/refresh-token`
- **Authorization**: `token_string`
- **Responses**:
  - **200**: gửi token mới cho người dùng và vô hiệu hóa token cũ.
  - **401**: Token không hợp lệ
  - **500**: Lỗi server

#### 7. Đăng xuất
- **Mô tả**: cho phép người dùng đăng xuất, backend sẽ vô hiệu hóa token.
- **Method**: POST
- **Endpoint**: `/api/v1/auth/logout`
- **Authorization**: `token_string`
- **Responses**:
  - **200**: token đã được vô hiệu hóa thành công
  - **400**: Token không được cung cấp

---
### II. Hiện thực API về admin management

#### 1. Lấy thông tin tổng quát về tất cả user
- **Mô tả**: Quản trị viên có thể truy xuất danh sách tất cả người dùng trong hệ thống. Trả về thông tin cơ bản về mỗi người dùng.
- **Method**: GET
- **Endpoint**: `/api/v1/admin/users`
- **Authorization**: `token_string`
- **Responses**:
  - **200**:  trả về list các thông tin của tất cả người dùng bao gồm id, email, username, vip_level và status.
  - **403**: Token không hợp lệ (không phải admin)
  - **500**: Lỗi server

#### 2. Lấy thông tin chi tiết của user

- **Mô tả**: Truy xuất thông tin chi tiết người dùng bằng cách cung cấp ID người dùng. Trả về thông tin người dùng như tên người dùng, email, cấp độ VIP, v.v.
- **Method**: GET
- **Endpoint**: `/api/v1/admin/user/{user_id}`
- **Authorization**: `token_string`
- **Header**: `User_id`
- **Responses**:
  - **200**: Trả về thông tin chi tiết của một user cụ thể thành công.
  - **400**: User Id không hợp lệ.
  - **403**: Token không hợp lệ (không phải admin).
  - **404**: Không tìm thấy user với id cung cấp.
  - **500**: Lỗi server.

#### 3. Lấy thông tin tổng quát về danh sách các lịch sử thanh toán

- **Mô tả**: Truy xuất lịch sử thanh toán cho tất cả người dùng. Quản trị viên có thể xem tất cả các khoản thanh toán một cách tổng quát.
- **Method**: GET
- **Endpoint**: `/api/v1/admin/payment-history`
- **Authorization**: `token_string`
- **Responses**:
  - **200**: Trả về danh sách tất cả các lịch sử thanh toán một cách tổng quát bao gồm `order_id`, `userID`, `order_info`, `amount`, và `transaction_status`.
  - **403**: Token không phải của admin.
  - **500**: Lỗi server.


#### 4. Lấy thông tin lịch sử thanh toán của một user cụ thể

- **Mô tả**: Truy xuất lịch sử thanh toán của một người dùng cụ thể dựa trên ID người dùng được cung cấp. Endpoint này chỉ giới hạn quyền truy cập cho quản trị viên và trả về danh sách các giao dịch thanh toán liên quan đến người dùng.
- **Method**: GET
- **Endpoint**: `/api/v1/admin/payment-history/{user_id}`
- **Authorization**: `token_string`
- **Header**: `User_id`
- **Responses**:
  - **200**: Trả về danh sách lịch sử thanh toán của một người dùng cụ thể thành công.
  - **400**: User id không hợp lệ hoặc thiếu `user_id`.
  - **403**: Token không hợp lệ (không phải admin).
  - **404**: Không tìm thấy user với id cung cấp.
  - **500**: Lỗi server.

#### 5. Cấm 1 tài khoản người dùng

- **Mô tả**: Cấm tài khoản người dùng bằng cách đặt trạng thái tài khoản thành không hoạt động. Quản trị viên có thể sử dụng điều này để cấm tài khoản người dùng.
- **Method**: PUT
- **Endpoint**: `/api/v1/admin/user/{user_id}/ban`
- **Authorization**: `token_string`
- **Header**: `User_id`
- **Responses**:
  - **200**: Cấm tài khoản người dùng thành công.
  - **400**: User id không hợp lệ hoặc thiếu `user_id`.
  - **403**: Token không hợp lệ (không phải admin).
  - **404**: Không tìm thấy user với id cung cấp.
  - **500**: Lỗi server.

#### 6. Kích hoạt lại tài khoản người dùng

- **Mô tả**: Kích hoạt tài khoản người dùng bằng cách đặt trạng thái tài khoản thành hoạt động. Quản trị viên có thể sử dụng chức năng này để kích hoạt lại tài khoản bị cấm.
- **Method**: PUT
- **Endpoint**: `/api/v1/admin/user/{user_id}/active`
- **Authorization**: `token_string`
- **Header**: `User_id`
- **Responses**:
  - **200**: Kích hoạt lại tài khoản người dùng thành công.
  - **400**: User id không hợp lệ hoặc thiếu `user_id`.
  - **403**: Token không hợp lệ (không phải admin).
  - **404**: Không tìm thấy user với id cung cấp.
  - **500**: Lỗi server.

#### 7. Xóa tài khoản người dùng


- **Mô tả**: Xóa một người dùng khỏi hệ thống bằng cách cung cấp ID người dùng. Người dùng sẽ bị xóa vĩnh viễn khỏi cơ sở dữ liệu.
- **Method**: DELETE
- **Endpoint**: `/api/v1/admin/user/{user_id}`
- **Authorization**: `token_string`
- **Header**: `User_id`
- **Responses**:
  - **200**: Xóa tài khoản người dùng thành công.
  - **400**: User id không hợp lệ hoặc thiếu `user_id`.
  - **403**: Token không hợp lệ (không phải admin).
  - **404**: Không tìm thấy user với id cung cấp.
  - **500**: Lỗi server.

---
