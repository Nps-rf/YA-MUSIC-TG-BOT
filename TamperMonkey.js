// ==UserScript==
// @name         Yandex Music Track Collector
// @namespace    http://tampermonkey.net/
// @version      0.1
// @description  Collect current track from Yandex Music
// @author       Nikolai Pikalov (https://github.com/Nps-rf)
// @match        https://music.yandex.ru/*
// @grant        GM_xmlhttpRequest
// ==/UserScript==

(function() {
    'use strict';

    function sendTrackInfoToServer(track) {
        GM_xmlhttpRequest({
            method: "POST",
            url: "http://localhost:8080/set-last-track",
            data: JSON.stringify(track),
            headers: {
                "Content-Type": "application/json"
            }
        });
    }

    function checkTrack() {
        const currentTrack = externalAPI.getCurrentTrack();

        if (currentTrack) {
            sendTrackInfoToServer(currentTrack);
        }
    }

    setInterval(checkTrack, 7000);
})();

