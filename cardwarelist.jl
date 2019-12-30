include("CardWare.jl")
using .CardWare: cardware_words, suits, values, symbols, cards

using Random
using ArgParse


printwords(wordlist::Vector) = begin
    let i = 1
        for c1 in 1:52
            s1 = Int(floor((c1-1)/13)) + 1
            suit1 = suits[s1]
            val1 = values[(c1-1)%13 + 1]
            card1 = cards[c1]
            for c2 in 1:52
                if c1 == c2
                    continue
                end
                s2 = Int(floor((c2-1)/13)) + 1
                suit2 = suits[s2]
                val2 = values[(c2-1)%13 + 1]
                card2 = cards[c2]

                word = wordlist[i]

                println("$card1$card2 ($val1$suit1+$val2$suit2) $word")

                i += 1

            end
        end
    end
end

printsymbols(rng::AbstractRNG) = begin
    s = shuffle(rng, symbols)
    println("   Blk Red")
    let i = 1
        for v in values
            println("$v   $(s[i])   $(s[i+1])")
            i += 2
        end
    end
end

printcaps(rng::AbstractRNG) = begin
    cap = rand(rng, ["Blk", "Red"])
    println("Capital letter: $cap")
end

printusage() = begin
    println("Usage: cardwarelist.jl [min_word_length] [seed]")
    println("min_word_length defaults to 4")
    println("if seed is not given, use entropy from the OS RNG to generate random numbers")
end

function parseargs()
    s = ArgParseSettings()
    @add_arg_table s begin
        "--min", "-m"
            help="minimum word length"
            arg_type=Int
            default=4
        "--seed", "-s"
            help="random seed (if no seed given, use entropy from the OS RNG)"
            arg_type=UInt32
            action=:append_arg
    end

    parse_args(s, as_symbols=true)
end

main() = begin
    parsed_args = parseargs()
    min_length = parsed_args[:min]
    rng = RandomDevice()
    if haskey(parsed_args, :seed) && !isempty(parsed_args[:seed])
        rng = MersenneTwister(parsed_args[:seed])
    end

    words = cardware_words(min_length=min_length, rng=rng)
    printwords(words)
    println()
    printsymbols(rng)
    println()
    printcaps(rng)
end

main()
