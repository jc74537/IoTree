
function addSong(name, id) {
    $("#songPicker").append(`<option id="newSongSel"></option>`);
    $("#newSongSel").text(name).attr("id", id).attr("value", name);
}
function addPattern(name, id) {
    $("#patternPicker").append(`<option id="newPatternSel"></option>`);
    $("#newPatternSel").text(name).attr("id", id).attr("value", name);
}
function populateSongs() {
    $.ajax({ //populate users dropdown
        url: "/api/song/list",
        type: 'GET',
        contentType: 'application/json',
        success: response => {
            console.log(response);
            songs = response;
            for (let song in songs) {
                addSong(songs[song].SongName, songs[song].SongID);
            }
        }
    });
}
function populatePatterns() {
    $.ajax({ //populate users dropdown
        url: "/api/lights/list",
        type: 'GET',
        contentType: 'application/json',
        success: response => {
            console.log(response);
            Patterns = response;
            for (let Pattern in Patterns) {
                addPattern(Patterns[Pattern].PatternName, Patterns[Pattern].PatternID);
            }
        }
    });
}
function playSong() {
    let option = {
        SongName: $("#songPicker").find("option:selected").val(),
        SongID: $("#songPicker").find("option:selected").attr("id")
    };
    if (option.name !== "Choose a song...") {
        console.log(option);
        $.ajax({ //populate users dropdown
            url: "/api/song/"+option.SongID,
            type: 'POST'
        });
        //location.reload()
        //$("#songPicker:first-child").attr("selected", true);
        //$("#songPicker").val.prop('selected', true);
    } else {
        alert("Pick a song.");
    }
}
function playPattern() {
    let option = {
        patternName: $("#patternPicker").find("option:selected").val(),
        patternID: $("#patternPicker").find("option:selected").attr("id")
    };
    if (option.name !== "Choose a pattern...") {
        console.log(option);
        $.ajax({ //populate users dropdown
            url: "/api/lights/"+option.patternID,
            type: 'POST'
        });
        //location.reload()
        //$("#patternPicker:first-child").attr("selected", true);
        //$("#patternPicker").val.prop('selected', true);
    } else {
        alert("Pick a pattern.");
    }
}
populateSongs();
populatePatterns();