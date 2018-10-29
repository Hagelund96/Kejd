package handler

/*
parts := strings.Split(r.URL.Path, "/")

switch len(parts) {
//handling /api/igc/ and /api/igc/id
case 5:
if parts[4] == "" {
replyWithAllTracksId(w, _struct.Db)
} else if checkId(parts[4]) {
replyWithTracksId(w, _struct.Db, parts[4])
} else {
http.Error(w, http.StatusText(404), 404)
}
case 6:
//handling /api/igc/id/ and /api/igc/id/field

if parts[5] == "" {
if !checkId(parts[4]) //!idExists
{
http.Error(w, "ID out of range.", http.StatusNotFound)
return
} else {
replyWithTracksId(w, _struct.Db, parts[4])
}
} else {
if checkId(parts[4]) {
replyWithField(w, _struct.Db, parts[4], parts[5])
} else {
http.Error(w, "Not a valid request", http.StatusBadRequest)
}
}
//handling /api/igc/id/field/
case 7:
if parts[6] == "" {
if checkId(parts[4]) {
replyWithField(w, _struct.Db, parts[4], parts[5])
} else {
http.Error(w, "Not a valid request", http.StatusBadRequest)
}
} else {
http.Error(w, "Not a valid request", http.StatusBadRequest)
}
}
*/