// ==UserScript==
// @name         Yandex Music Track Collector
// @namespace    http://tampermonkey.net/
// @version      0.3
// @description  Collect current track from Yandex Music
// @author       Nikolai Pikalov (https://github.com/Nps-rf)
// @match        https://music.yandex.ru/*
// @grant        GM_xmlhttpRequest
// ==/UserScript==

(function () {
    'use strict';
    const TIMEOUT = 5000;  // 5 seconds
    const PROTOCOL = "https://";
    const BACKEND_URL = "http://localhost:8080";
    const ENDPOINT = "/set-last-track";

    const method = "POST";
    const headers = {"Content-Type": "application/json"}

    const url = `${PROTOCOL}${BACKEND_URL}${ENDPOINT}`

    const sendTrackInfoToServer = data => GM_xmlhttpRequest({
        method,
        url,
        data,
        headers,
    });


    const checkTrack = () => {
        const currentTrack = externalAPI.getCurrentTrack();
        (currentTrack && externalAPI.isPlaying()) && sendTrackInfoToServer(JSON.stringify(currentTrack));
    }

    setInterval(checkTrack, TIMEOUT);
})();

