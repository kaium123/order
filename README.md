### Project Setup
### Clone the repository
```bash
git clone https://github.com/kaium123/order.git
```
### Install Dependecies
```bash
go mod tidy
```
#### Prerequisites:
- Docker (for running containers)
- Docker Compose (for managing multi-container applications)
#### 1. **Running the Application Locally**
#### Command:
Run the following command to start the **application** locally, **database**, **Redis**, and **Consul** in docker container:

```bash
./run_local.sh
```

#### 2. **Running the Application with Docker**
Run the following command to run the **database**, **Redis**, **Consul**, and the **application** within Docker containers:
```bash
./run_docker.sh
```

### server port and configution stored in config.yaml and config.docker.yaml

### Orders API
#### Features
#### 1. **Login**
   - **Endpoint**: `/api/v1/login`
   - **Description**:  This endpoint allows users to authenticate with their credentials. Upon successful login, the API returns an access token and a refresh token that can be used to interact with the other API endpoints.
   - **Input**:  
     ```json
     {
         "username": "01901901901@mailinator.com",
         "password": "321dsa"
     }
     ```
   - **Response**:  
     ```json
     {
         "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzE4Njg4NTMsInVzZXJfaWQiOjF9.I0ePj4bvQsCW5rODb4uPBjSRt1bUCVmDIMcBYurhEKQ",
         "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzI0NzAwNTMsInVzZXJfaWQiOjF9.LuQvDKekCPZGIxI9KlR0NxVJcUSnz6P6YQUppuj9EgQ",
         "token_type": "Bearer",
         "expire_in": 1731868853
     }
     ```

#### 2. **Logout**
   - **Endpoint**: `/api/v1/logout`
   - **Description**:  
     Allows the user to invalidate their session. The provided access and refresh tokens are revoked, ensuring the user is logged out securely.
   - **Input**:  
     - Access Token (in the header)  
   - **Response**:  
     ```json
     {
         "message": "Successfully logged out",
         "code": "200",
         "type": "success"
     }
     ```

#### 3. **Create Order**
   - **Endpoint**: `/api/v1/orders`
   - **Description**:  Enables users to place new orders by providing necessary details.
   - **Input**:  
     ``` json
        {
            "store_id": 131172,
            "merchant_order_id": "123",
            "recipient_name": "kaium",
            "recipient_phone": "01875113838",
            "recipient_address": "banani, gulshan 2, dhaka, bangladesh",
            "recipient_city": 1,
            "recipient_zone": 1,
            "recipient_area": 1,
            "delivery_type": 2,  ///1 - pickup, 2 - delivery, 
            "item_type": 2,  /// 1, document, 2 - parcel, 3 - other
            "special_instruction": "please provide as soon as possible",
            "item_quantity": 1,
            "item_weight": 0.5,
            "amount_to_collect": 12000,
            "item_description": "this is description"
        }
        ```
   - **Response**:  
     ```json
     {
         "message": "Order Created Successfully",
         "code": "200",
         "type": "success",
         "data": {
             "consignment_id": "DA241117FABNUY",
             "merchant_order_id": "123",
             "order_status": "Pending",
             "delivery_fee": 60
         }
     }
     ```

#### 4. **Cancel Order**
   - **Endpoint**: `/api/v1/orders/{CONSIGNMENT_ID}/cancel`
   - **Description**:  Users can cancel an existing order by providing the order ID. Orders that are already processed or delivered cannot be canceled.
   - **Input**:  CONSIGNMENT_ID (in the path)  
   - **Response**:  
     ```json
     {
         "message": "Order Cancelled Successfully",
         "code": "200",
         "type": "success"
     }
     ```

#### 5. **Fetch Order List**
   - **Endpoint**: `api/v1/orders/all`
   - **Description**:  Retrieve a list of all orders placed by the user. Supports filters.
   - **Input**:  
     - Filters: `?limit=1&page=2&transfer_status=1&archive=0`  
   - **Response**:  
     ```json
     {
         "message": "Orders successfully fetched.",
         "code": "200",
         "type": "success",
         "data": {
             "orders": [
                 {
                     "order_consignment_id": "DA241117SIHWXX",
                     "order_created_at": "2024-11-17T17:34:07.032741Z",
                     "order_description": "this is description",
                     "merchant_order_id": "123",
                     "recipient_name": "kaium",
                     "recipient_address": "banani, gulshan 2, dhaka, bangladesh",
                     "recipient_phone": "01875113838",
                     "order_amount": 12000,
                     "total_fee": 180,
                     "instruction": "please provide as soon as possible",
                     "order_type_id": 1,
                     "cod_fee": 120,
                     "promo_discount": 0,
                     "discount": 0,
                     "delivery_fee": 60,
                     "order_status": "Pending",
                     "order_type": "Delivery",
                     "item_type": "Parcel"
                 }
             ],
             "total": 17,
             "current_page": 2,
             "per_page": 1,
             "total_in_page": 1,
             "last_page": 17
         }
     }
     ```



### 3. **Optimizations**

#### Singleton Design Pattern for Database and Redis
To ensure that there is only **one instance** of the database and Redis connection throughout the application, I have utilized the **Singleton Design Pattern**. This pattern ensures that the database and Redis connections are created once, and that instance is reused whenever needed.

#### Using Indexing for Faster Data Retrieval
To improve the speed of data retrieval from the database, **indexes** are created on frequently queried columns. Indexing is particularly beneficial when working with large datasets and performing read-heavy operations. By adding indexes, the database can quickly locate rows without needing to scan the entire table.

#### Caching with Redis
To enhance performance, especially when dealing with read-heavy data, we have implemented **Redis** as a caching layer. Redis stores frequently accessed data in memory, allowing for fast retrieval without querying the database every time. 

#### Object-Oriented Programming (OOP) Concepts
I follow **Object-Oriented Programming (OOP)** principles to make the application maintainable, scalable, and reusable:

#### Postman Collection :
https://api.postman.com/collections/32612509-05d4b762-1ef6-4191-ac8e-3d3ecabf9b0d?access_key=PMAT-01JCXKFNJVG1GABDAWEMN1VHQ8
