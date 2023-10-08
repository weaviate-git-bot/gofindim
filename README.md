# Gofindim

`gofindim` is an experimental, exploratory project of mine to leverage [CLIP models](https://openai.com/research/clip) to enable local image search using both text and image queries. It is written in Go for educational purposes, but the better practical choice would probably NodeJS for a more cohesive tech stack.

#### Features:
* CLI commands for individual files
* Uses [Go-chi](https://github.com/go-chi/chi) as a backend webserver
* Uses [Weaviate](https://github.com/weaviate/weaviate) to store vectors
    * Differing from typical Weaviate usage is that images are **not** stored in the database, only their paths and vectors
* ReactJS frontend UI
* Supports scanning directories and individual images.
* Uses [air](https://github.com/cosmtrek/air) for development

#### Limitations and Quirks:
* Only works on Linux
* Assumes `/web` is **not** a valid directory
* Since only image paths and vectors are stored in the database, once a file is scanned,  it assumes the file will not change locations.
    * To address this, a wrapper for `mv` is planned that will modify the paths of matching images in the database

### Installation & Usage

#### 1. Start the database
1. In the repo directory run:
  * `mkdir .weaviate_data`
  * `docker compose up -d`
#### 2. Build the frontend files
1. `cd web/frontend`
2. `yarn install`
3. `yarn run build`
#### 3. Build the backend
1. `go build -o gofindim`
2. Initialize the database
    * `./gofindim add -G` will create the Image class in the database

#### 4. Basic Usage
1. `./gofindim web` will start a webserver at `localhost:8888`
    * Navigate to http://localhost:8888/web/ to reach the web UI
    * While looking for similar images
        * Clicking an image will apply that image to the query
        * Modify any needed parameters and hit apply to apply changes
        * The submit button will start a new query using only the text with no image applied
        
2. `./gofindim add` will scan individual images into the database
    * Use the `-A` flag to scan directory contents
3.  `./gofindim search` allows for reverse-image search if provided a path, otherwise it will perform a text search

#### 5. For development:
1. Use `air web` to start the webserver
2. Use `yarn run start`  inside `web/frontend` that will start on http://localhost:4000

### Images:
Web UI Landing
![Landing](https://res.cloudinary.com/dbtfu4e73/image/upload/v1696734870/main_b36bd4.png)
Browse UI
![Browse UI](https://res.cloudinary.com/dbtfu4e73/image/upload/v1696734870/browse_n2kmlu.png)
Searching using text
![Text search](https://res.cloudinary.com/dbtfu4e73/image/upload/v1696734871/similar1_naqur5.png)
Searching using both text and image
![Text and Image search](https://res.cloudinary.com/dbtfu4e73/image/upload/v1696734871/similar2_kmyfpn.png)


