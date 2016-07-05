package hash

// fork from lua
func HashString(str string) uint32 {
    byt := []byte(str)
    l := len(str)
    h := l
    step := (l >> 5) + 1
    for i := l; i >= step; i -= step {
        h = h ^ ((h << 5) + (h >> 2) + int(byt[i-1]))
    }
    if h == 0 {
        return uint32(1)
    }
    return uint32(h)
}
