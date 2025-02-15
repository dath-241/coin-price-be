# **Coin-price-be-go – GVHD: Thầy Lê Đình Thuận**
- [Giới thiệu dự án](#1-giới-thiệu-dự-án)
- [Giới thiệu nhóm](#2-giới-thiệu-nhóm)
- [Công nghệ sử dụng](#3-công-nghệ-sử-dụng)
- [Tài liệu về System Requirement Specification (SRS)](#4-tài-liệu-về-system-requirement-specification-srs)
- [Các chức năng chính](#5-các-chức-năng-chính-api-doc)
    - [Authentication](#51-authentication)
    - [Admin](#52-admin)
    - [User](#53-user)
    - [Coin price](#54-coin-price)
    - [Trigger](#55-trigger)
    - [Payment](#56-payment)
- [Cấu trúc các nhánh trên hệ thống](#6-cấu-trúc-các-nhánh-trên-hệ-thống)
    - [Cấu trúc nhánh main](#61-cấu-trúc-nhánh-main)
    - [Cấu trúc nhánh develop](#62-cấu-trúc-nhánh-develop)
    - [Cấu trúc nhánh production](#63-cấu-trúc-nhánh-production)
- [Tham thảo thêm](#7-tham-khảo-thêm)
## **1. Giới thiệu dự án**
**Coin-price-be-go** là hệ thống backend với chức năng chính là hỗ trợ người dùng việc theo dõi thông tin thị trường tiền điện tử (giá spot, future, kline,...), tạo các trigger thông báo biến động về giá sử dụng các API từ Binance, CoinMarketCap, CoinGecko.
## **2. Giới thiệu nhóm**
**Tên nhóm**: Nhóm Ngọc Hiền
**Danh sách thành viên**: 7 thành viên
|    MSSV   |         Họ tên         |  Role  | Github |Nhiệm vụ|
|:---------:|:----------------------:|:------:|:------:|:------:|
|  2211024  |      Lê Ngọc Hiền      | PO     |[HienLe2004](https://github.com/HienLe2004), [LeHien6601](https://github.com/LeHien6601)|Tổ chức họp, đánh giá tiến độ, xây dựng cấu trúc dự án, phân chia công việc, hiện thực CI/CD|
|  2211101  |   Nguyễn Thanh Hoàng   | Dev    |[NTHKiris](https://github.com/NTHKiris)|Hiện thực các API về alerts và viết unit test, API doc kèm theo|
|  2210768  |    Nguyễn Văn Đoàn     | Dev    |[DoanJackson](https://github.com/DoanJackson), [anhchienne](https://github.com/anhchienne)|Hiện thực các API về coin-price và viết unit test, API doc kèm theo|
|  2211144  | Nguyễn Trịnh Ngọc Huân | Dev    |[huannguyen2114](https://github.com/huannguyen2114)|Hiện thực các API về coin-price và viết unit test, API doc kèm theo|
|  2212962  |     Trần Quang Tác     | Dev    |[tacsquang](https://github.com/tacsquang)|Hiện thực các API về admin, auth, user, payment và viết unit test, API doc kèm theo|
|  2210871  |     Quách Khải Hào     | Dev    |[quachkhaihao](https://github.com/quachkhaihao)|Hiện thực các API về alerts và viết unit test, API doc kèm theo|
|  2212922  |    Nguyễn Quang Sáng   | Dev    |[Sangquangnqs](https://github.com/Sangquangnqs), [millerbright](https://github.com/millerbright)|Hiện thực các API về admin, auth, user, payment và viết unit test, API doc kèm theo|

## **3. Công nghệ sử dụng**
| Công nghệ | Mô tả |
|-------------|----------|
|[Golang](https://go.dev/)|Ngôn ngữ phát triển hệ thống backend với hiệu năng cao, hỗ trợ tốt cho xây dựng hệ thống phân tán và microservices|
|[Gin](https://gin-gonic.com/)|Framework của Go, giúp xây dựng các API RESTful nhẹ nhàng và nhanh chóng|
|[Swaggo](https://github.com/swaggo/swag)|Tạo tài liệu cho API tự động và hỗ trợ tương tác giúp cho developer dễ dàng hiểu để sử dụng API|
|[Git](https://git-scm.com/), [GitHub](https://github.com/)|Quản lý phiên bản code, hỗ trợ các chức năng cho cộng tác như pull request, issue tracking, CI/CD|
|[Docker](https://www.docker.com/), [DockerHub](https://hub.docker.com/), [DockerCompose](https://docs.docker.com/compose/)|Công cụ container hóa ứng dụng, xây dựng, lưu trữ các image, định nghĩa và khởi chạy đa container|
|[Mailjet](https://www.mailjet.com/)|Dịch vụ hỗ trợ gửi email xác nhận, thông báo đến người dùng|
|[MongoDB](https://www.mongodb.com/)|Cơ sở dữ liệu NoSQL linh hoạt cho việc lưu trữ thông tin người dùng, các trigger, lịch sử thanh toán,...|
|[MomoAPI](https://developers.momo.vn/v2/#/)|API của ví điện tử MOMO để tích hợp tính năng thanh toán để nâng cấp VIP vào hệ thống|
|[BinanceAPI](https://www.binance.com/en/binance-api), [CoinGecko](https://www.coingecko.com/en/api), [CoinMarketCap](https://coinmarketcap.com/api/)|API hỗ trợ lấy dữ liệu thị trường giao dịch tiên ảo realtime|
|[Testify](https://github.com/stretchr/testify)|Framework testing của Go hỗ trợ viết unit tests đơn giản và hiệu quả|
## **4. Tài liệu về System Requirement Specification (SRS)**
Sử dụng chung SRS với nhóm coin-price-web: [link](https://github.com/dath-241/coin-price-web/blob/main/Documents/System%20Requirement%20Specification/Crypto_Market_SRS.pdf)
## **5. Các chức năng chính** 
Link SwaggerUI cho API doc [(link)](https://a1-price.thuanle.me/docs/index.html)

Các demo chạy chức năng [auth, admin, user, payment](https://drive.google.com/drive/u/0/folders/1K-4gh6WLLL45MHfxtsAJNu-4GHYpwoAY), demo [trigger](https://github.com/dath-241/coin-price-be-go/issues/4), demo [coin_price](https://documenter.getpostman.com/view/40206908/2sAYBaAVZH)
### **5.1. Authentication**
Cung cấp các chức năng chính cho người dùng về xác thực như: đăng ký, đăng nhập, đăng nhập bằng google, quên mật khẩu, đổi mật khẩu, đăng xuất, làm mới token.
[Chi tiết](https://github.com/dath-241/coin-price-be-go/blob/main/docs/auth_admin_user_payment.md#i-hi%E1%BB%87n-th%E1%BB%B1c-api-v%E1%BB%81-x%C3%A1c-th%E1%BB%B1c-authentication)
### **5.2. Admin**
Cung cấp các chức năng chính cho admin như: lấy danh sách người dùng, lấy thông tin người dùng, xóa người dùng, cấm người dùng, bỏ cấm người dùng, lấy lịch sử thanh toán của người dùng và của hệ thống.
[Chi tiết](https://github.com/dath-241/coin-price-be-go/blob/main/docs/auth_admin_user_payment.md#ii-hi%E1%BB%87n-th%E1%BB%B1c-api-v%E1%BB%81-admin-management)
### **5.3. User**
Cung cấp các chức năng chính cho người dùng liên quan đến dữ liệu người dùng như: lấy thống tin người dùng, cập nhật thông tin người dùng, xóa tài khoản người dùng, đổi mật khẩu người dùng, đổi email người dùng, lấy lịch sử thanh toán, lấy danh sách các trigger của người dùng, gửi thông báo trigger đến email người dùng.
[Chi tiết](https://github.com/dath-241/coin-price-be-go/blob/main/docs/auth_admin_user_payment.md#iii-hi%E1%BB%87n-th%E1%BB%B1c-api-v%E1%BB%81-user-management)
### **5.4. Coin price**
Cung cấp các chức năng chính cho người dùng lấy dữ liệu thị trường tiền ảo: lấy giá spot realtime, lấy giá future realtime, lấy funding rate realtime, lấy kline, lấy marketcap,...
[Chi tiết](https://github.com/dath-241/coin-price-be-go/blob/main/docs/coin_price.md)
### **5.5. Trigger**
Cung cấp các chức năng chính cho người dùng về thiết lập các trigger: tạo trigger, lấy trigger, xóa trigger, lấy danh sách new/delisted symbol, tạo trigger cho new/delisted symbol, khởi chạy trigger, dừng trigger, tạo trigger cho advanced indicators,...
[Chi tiết](https://github.com/dath-241/coin-price-be-go/blob/main/docs/trigger.md)
### **5.6. Payment**
Cung cấp các chức năng chính cho người dùng về thanh toán: khởi tạo thanh toán nâng cấp VIP qua Momo, gọi thanh toán qua Momo, kiểm tra trạng thái thanh toán qua Momo, xác nhận thanh toán và nâng cấp VIP qua Momo.
[Chi tiết](https://github.com/dath-241/coin-price-be-go/blob/main/docs/auth_admin_user_payment.md#iv-hi%E1%BB%87n-th%E1%BB%B1c-api-v%E1%BB%81-payment)
## **6. Cấu trúc các nhánh trên hệ thống**
### **6.1. Cấu trúc nhánh** [main](https://github.com/dath-241/coin-price-be-go/tree/main)
```plaintext
.
├── /reports # Các meeing minutues và đánh giá hàng tuần của nhóm
├── /docs # Các mô tả chi tiết về các chức năng chính của hệ thống
├── /backend # Mã nguồn chính của dự án
└── README.md # File mô tả chính về hệ thống
```
### **6.2. Cấu trúc nhánh** [develop](https://github.com/dath-241/coin-price-be-go/tree/develop)
```plaintext
.
├── /.github
│   └── /workflows 
│       └── dev.yml # CI/CD cho nhánh develop
├── /backend # Mã nguồn chính của dự án
└── README.md # File mô tả phân công công việc
```
### **6.3. Cấu trúc nhánh** [production](https://github.com/dath-241/coin-price-be-go/tree/production)
```plaintext
.
├── /.github
│   └── /workflows 
│       └── production.yml # CI/CD cho nhánh production
└── docker-compose.yml # File cấu hình docker-compose cho VPS khi chạy 
```


## **7. Tham khảo thêm**
Các [issue](https://github.com/dath-241/coin-price-be-go/issues)

Các demo [auth, admin, user, payment](https://drive.google.com/drive/u/0/folders/1K-4gh6WLLL45MHfxtsAJNu-4GHYpwoAY), demo [trigger](https://github.com/dath-241/coin-price-be-go/issues/4), demo [coin_price](https://documenter.getpostman.com/view/40206908/2sAYBaAVZH)

