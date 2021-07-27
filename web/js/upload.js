(async function () {
    main();
}());

function main() {
    const  BASE_URL = "http://localhost:8000";
    // const BASE_URL = "https://evening-atoll-24533.herokuapp.com";
    let resultsElement = document.getElementById("results");
    let uploadFormElement = document.getElementById("uploadForm");
    let fileInputElement = document.getElementById("fileInput");
    let eventIdElement = document.getElementById("eventId");
    uploadFormElement.addEventListener('submit', async (e) => {
        e.preventDefault();

        let requestOptions = {
            method: 'GET',
            redirect: 'follow'
        };

        try {
            let formData = new FormData();
            formData.append("file", fileInputElement.files[0], "/C:/Users/mikha/Downloads/text.zip");
            formData.append("eventId", "1");

            let requestOptions = {
                method: 'POST',
                body: formData,
                redirect: 'follow'
            };


            const response = await fetch(BASE_URL + "/api/v1/picture/detectText", requestOptions);
            const data = await response.text();
            console.log(data);
        } catch (err) {
            console.log('error', err);
        }
    });
}