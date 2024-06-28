function loadPage(page) {
    console.log('Loading page:', page);
    const contentDiv = document.getElementById('content');
    const xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {
            contentDiv.innerHTML = xhr.responseText;
        }
    };
    xhr.open('GET', `pages/${page}`, true);
    xhr.send();
}