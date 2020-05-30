# SellerApp_Interview_Assignment
## Simple REST Services implemented in Golang

## Description

There are 2 POST API's in this project . One API (api-1) is for scraping product information form the url sent in the request body,After
getting the product information it is marshalled into json and sent to the other API (api-2) in a request body .
The second API gets the product info and inserts/updates in a document store (mongoDB) with other info like creation and update
time .


## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

1. Clone this repository.
```
https://github.com/zirin12/SellerApp_Interview_Assignment.git
```

2. Build and launch using docker.
  Make sure Docker is running First. 
  ```
  sudo service docker status
  ```
  If it isn't running , Then start it 
  
  ```
  sudo service docker start
  ```

  Run the following commands in the project directory.
  ```
  docker-compose up -d --build
  ```
  This will create 3 docker containers which are api-1 , api-2 and the mongoDB container . 

  To know more information about these containers like it's names , Ports exposed and status:
  ```
  docker ps
  ```

3. To Test the service , you can use POSTMAN or CURL or any api testing tool

  First we have to give a product url to api-1 in the form of a post request 
  Send a POST request to to http://localhost:9080/product with url as json in request body.
  Sample CURL call is below : 
  ```
  curl --location --request POST 'localhost:9080/product' --header 'Content-Type: application/json' --data-raw 
  '{
    "url" : "https://www.amazon.com/3003348-PlayStation-4/dp/B071CV8CG2/"
  }'
  ```
  It will return a response which contains the scraped information of the product :
  ```
  {
        "url": "https://www.amazon.com/3003348-PlayStation-4/dp/B071CV8CG2/",
        "product": {
            "title": "PlayStation 4 Slim 1TB Console",
            "imageURL": "https://images-na.ssl-images-amazon.com/images/I/71PGvPXpk5L._AC_SX425_.jpg",
            "description": "Incredible games; Endless entertainment",
            "rating": "4.7 out of 5 stars",
            "price": "$384.00",
            "totalReviews": "4,734 ratings"
        }
  }
  ```
  The above API internally call the other API as explained before which persists it to a database . To see all the items currently
  in mongoDB , we can hit a GET API from the second service to get all the records in database . The GET API call is to the url
  http://localhost:8080/products . Sample CURL call below :
   ```
   curl --location --request GET 'localhost:8080/products' --header 'Content-Type: application/json'
   ```
  It will give a response like the one below of all the records :
  ```
  [
      {
          "id": "5ed24997b922eb16b6b9cecf",
          "url": "https://www.amazon.com/3003348-PlayStation-4/dp/B071CV8CG2/",
          "created_at": "2020-05-30 11:55:02.986961202 +0000 UTC",
          "product": {
              "title": "PlayStation 4 Slim 1TB Console",
              "imageURL": "https://images-na.ssl-images-amazon.com/images/I/71PGvPXpk5L._AC_SX425_.jpg",
              "description": "Incredible games; Endless entertainment",
              "rating": "4.7 out of 5 stars",
              "price": "$384.00",
              "totalReviews": "4,734 ratings"
          },
          "updated_at": "2020-05-30 11:55:02.986961202 +0000 UTC"
      },
      {
          "id": "5ed24b4fb922eb16b6b9cf02",
          "url": "https://www.amazon.com/Wyze-Indoor-Wireless-Detection-Assistant/dp/B076H3SRXG/",
          "created_at": "2020-05-30 12:02:23.049956781 +0000 UTC",
          "product": {
              "title": "Wyze Cam 1080p HD Indoor Wireless Smart Home Camera with Night Vision, 2-Way Audio, Works with Alexa & the Google Assistant, One Pack, White - WYZEC2",
              "imageURL": "https://images-na.ssl-images-amazon.com/images/I/61B04f0ALWL._AC_SY355_.jpg",
              "description": "Live Stream from Anywhere in 1080p -1080p Full HD live streaming lets you see inside your home from anywhere in real time using your mobile device. While live streaming, use two-way audio to speak with your friends and family through the Wyze app.",
              "rating": "4.3 out of 5 stars",
              "price": "$25.98",
              "totalReviews": "39,615 ratings"
          },
          "updated_at": "2020-05-30 12:02:23.049956781 +0000 UTC"
      }
  ]
  ```

4. To shut down all the services:
  ```
  docker-compose down
  ```
5. If you want to access logs of each service then you can use the following command : 
  ```
  docker logs < Container name > -f --tail=50 # tail=50 gives the last 50 lines of the log
  ````
  Container name can be gotten from executing docker-compose ps


