package v1

//go:generate flatc --go -o generated ocap.fbs
//go:generate sh -c "sed -i 's/^package .*/package generated/' generated/*.go"
