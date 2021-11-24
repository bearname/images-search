(async function () {
    await main();
}());

async function main() {
    const  BASE_URL = "http://localhost:8000";
    // const BASE_URL = "https://thawing-sea-83431.herokuapp.com";
    let resultsElement = document.getElementById("results");
    let searchFormElement = document.getElementById("searchForm");
    let searchInputElement = document.getElementById("searchInput");
    let confidenceElement = document.getElementById("confidence");
    let nextElement = document.getElementById("next");
    const CONFIDENCE_DEFAULT = 85;
    const page = {limit: 20, offset: 0}
    const data = await searchImages(BASE_URL, {confidence: CONFIDENCE_DEFAULT, limit: page.limit, offset: page.offset});
    drawImages(data, resultsElement);
    searchFormElement.addEventListener('submit', async (e) => {
        e.preventDefault();
        try {
            let confidence = confidenceElement.value;
            let q = searchInputElement.value;
            const data = await searchImages(BASE_URL, {
                q: q,
                confidence: confidence,
                limit: page.limit,
                offset: page.offset
            });
            drawImages(data, resultsElement);
            console.log(data);
        } catch (err) {
            resultsElement.innerHTML = 'Not found';

            notFoundImages(resultsElement);
            console.log('error', err);
        }
    });
    nextElement.addEventListener('click', async (e) => {
        // let currentPage = Number.parseInt( nextElement.getAttribute('data-page'));
        e.preventDefault();
        page.offset += page.limit;
        try {
            let confidence = confidenceElement.value;
            let q = searchInputElement.value;
            const data = await searchImages(BASE_URL, {
                q: q,
                confidence: confidence,
                limit: page.limit,
                offset: page.offset
            });
            drawImages(data, resultsElement, false);
            console.log(data);
        } catch (err) {
            notFoundImages(resultsElement);
            console.log('error', err);
        }
    });
}

async function searchImages(baseUrl, options = {}) {
    try {
        let url = `${baseUrl}/api/v1/picture/search?`;
        if (options.hasOwnProperty('q')) {
            url += `number=${options.q}`;
        }
        if (options.hasOwnProperty('confidence')) {
            url += `&confidence=${options.confidence}`;
        }
        if (options.hasOwnProperty('eventId')) {
            url += `&eventId=${options.eventId}`;
        }
        if (options.hasOwnProperty('limit')) {
            url += `&limit=${options.limit}`;
        }
        if (options.hasOwnProperty('offset')) {
            url += `&offset=${options.offset}`;
        }

        const response = await fetch(url);
        if (!response.ok) {
            return [];
        }
        return await response.json()
    } catch (err) {
        console.log('error', err);
        return [];
    }
}

function notFoundImages(parentElement) {
    document.getElementById('alert').innerHTML = 'Not found';
}

function drawImages(data, parentElement, needClear = true) {
    const pictures = data.Pictures;
    if (pictures.length === 0) {
        if (needClear) {
            parentElement.innerHTML = '';
        }
        notFoundImages();
    } else {
        if (needClear) {
            parentElement.innerHTML = '';
        }
        let innerHtml = '';
        for (const picture of pictures) {
            innerHtml += drawImage(picture, innerHtml);
        }
        if (needClear) {
            console.log('needclear')
            parentElement.innerHTML = innerHtml;
        } else {
            console.log('notneedclear')
            let element = document.createElement('div');
            element.innerHTML = innerHtml
            parentElement.appendChild(element);
        }
    }
}

function drawImage(picture) {
    let picturePath = picture.Path;

    picturePath = (picturePath.substring(0, "https://".length) === 'https://') ? picturePath : "/images/" + picturePath;
    console.log(picturePath.substring(0, "https://".length));
    let innerHtml = `<li data-v-519eea8b="" class="loaded">
                         <a data-v-519eea8b data-lightbox="lightbox" class="">
                         <img src="${picturePath}" alt="${picture.PictureId}">
                        </a>`;
    innerHtml += `<span>`;
    for (const item of picture.TextDetections) {
        innerHtml += `#${item.DetectedText} ${item.Confidence}% `
    }

    ` </span> </li>`;
    return innerHtml;
}
