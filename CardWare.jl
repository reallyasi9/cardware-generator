module CardWare

export getwords, suits, values, symbols, cards

using HTTP: get
using MbedTLS: SSLConfig
using Random

const suits = split("♠♡♢♣", "")
const values = split("A23456789TJQK", "")
const _ace_of_spades = 0x1f0a1
const cards = [Char(_ace_of_spades + s * 16 + v) for s in 0:3 for v in 0:12]
const symbols = split(raw"_!#$%&()*+,-./:;<=>?@[]^{}", "")
const _words_url = "https://github.com/first20hours/google-10000-english/raw/master/google-10000-english-no-swears.txt"

function cardware_words(;min_length::Integer = 4, rng::AbstractRNG = RandomDevice())
    req = get(_words_url, sslconfig=SSLConfig(false))
    words = String(req.body)

    wordlist = split(words, "\n")

    acceptword(x) = !isempty(x) && !occursin(r"\s", x) && length(x) >= min_length

    filter!(acceptword, wordlist)

    # Random subset
    ndraws = 52*51
    shuffle!(rng, wordlist)
    wordsublist = wordlist[1:ndraws]

    sort!(wordsublist)
end

end
