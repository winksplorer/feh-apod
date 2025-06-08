# feh-apod

A simple Go program that uses `feh --bg-fill` to set your wallpaper to the current NASA APOD (Astronomy Picture of the Day)

## Building

Have Go and feh installed, and run `make`.

## Running

Have feh installed. Run the compiled binary called `feh-apod`.

## Running (with cron)

NASA's APOD usually updates at 00:00 ET. Convert that into your timezone, and maybe add 5 minutes of headroom.

For the path, simply just specify the full path to the compiled binary.

For example, I use this task on my personal system, to run `/home/wink/repos/feh-apod/feh-apod` at 21:05 PT.

```
5 21 * * * /home/wink/repos/feh-apod/feh-apod
```