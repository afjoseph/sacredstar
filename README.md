# SacredStar

This project implements an easy interface to conduct astrological research.

It can be used to easily cast a chart and:
- Run house calculations
- Do birthtime rectification
- Calculate different zodiacal positions
- Supports Vedic varga charts (up to D60)
- Calculates Vimshottari dashas
- Supports traditional and modern planets
- Supports sidereal and tropical charts

## Installation

Add the library to your Go project

    go get github.com/afjoseph/sacredstar

This library uses native code, so you'll need to run `CGO=1` when building

Also, you'll need to set up [SwissEph library]([swisseph](https://www.astro.com/swisseph/swephinfo_e.htm)). See below

### Setting up Swisseph

This project uses [Just task runner](https://github.com/casey/just) to run project commands. You can install it on an OSX machine with `brew install just`. See the project for other installation instructions.

This project uses [swisseph](https://www.astro.com/swisseph/swephinfo_e.htm) database, libraries and header files for calculation of planet and star positions.

All the required files are bundled in `./ext/swisseph/`.

In most cases, running `just build-and-symlink-swisseph` should be enough to set this project to work on your development machine.

#### Explanation

This is the directory breakdown of swisseph

- `./ext/swisseph`
    - the swisseph library from astro.com
    - `./original_project` contains the original swisseph source files, downloaded verbatim from https://www.astro.com/ftp/swisseph/src/
        - crudely with `wget -r --no-parent https://www.astro.com/ftp/swisseph/src/`
    - `./ephe` contains the ephemeris data files (i.e., position of stars and planets)
        - Here's a breakdown of each <https://www.saravali.de/articles/calc_ephemfiles.html>
        - Those files are used as is by the project. Don't remove or modify
    - `./include` contains the used header files
        - Those are copy-pasted from ./original_project
    - `./lib`
        - Contains the generated `libswe.a` file
        - `libswe.a` was generated with `just build-and-symlink-swisseph`

We use [pkg-config](https://www.freedesktop.org/wiki/Software/pkg-config/) to inform [cgo](https://pkg.go.dev/cmd/cgo) (the FFI responsible for calling C code from Go) where the swisseph libraries and header files are located. We're doing this using pkg-config, as opposed to a simple relative path, to maintain uniformity between building on a developer machine (e.g., using OSX) and a production instance (e.g., running a linux distro or using Docker with a linux distro).

If you look at the header of `./backend/swisseph/chart/chart.go`, you'll find this:

```
package chart

import (
	// #cgo pkg-config: swisseph
	// #include <stdio.h>
	// #include <errno.h>
	// #include "swephexp.h"
	"C"
)
```

This just means we're telling cgo to use `pkg-config --libs --cflags` during compilation. The `swisseph.pc` file is located in `./backend/ext/swisseph/`. You'll need to symlink it to your machine before backend development can start. You can do that with `just build-and-symlink-swisseph`. The commands were configured to work successfully on an OSX machine, but you might need to tweak it for your system.

Another reason we're doing this is because `libswe.a` (i.e., the main swisseph library) will use a different architecture depending on the machine it was built with (e.g., arm64 on an OSX machine with M1+ chips, and x86_64 on a production instance running a linux machine).

#### Setting up SwissEph on a Remote Machine

After building your project, you'd wanna deploy it somewhere. SwissEph needs to be symlinked properly on your remote machine.

For Ansible deployments, do something like this to copy, build and symlink swisseph on your remote machine

    - name: Copy swisseph directory to the remote machine
      copy:
        src: sacredstar/ext/swisseph
        dest: /home/myuser/app
        mode: '0755'

    - name: Set up Swisseph libraries and links
      vars:
        # Remember to
        swisseph_path: /home/myuser/app/swisseph
      shell: |
        rm -rf /usr/local/lib/pkgconfig/swisseph.pc
        rm -rf /usr/local/include/swisseph
        rm -rf /usr/local/lib/swisseph
        (cd {{ swisseph_path }}/original_project && make clean libswe.a)
        cp {{ swisseph_path }}/original_project/libswe.a {{ swisseph_path }}/lib/libswe.a
        mkdir -p /usr/local/lib/pkgconfig
        ln -sf {{ swisseph_path }}/swisseph.pc /usr/local/lib/pkgconfig/swisseph.pc
        ln -sf {{ swisseph_path }}/include /usr/local/include/swisseph
        ln -sf {{ swisseph_path }}/lib /usr/local/lib/swisseph
      args:
        executable: /bin/bash

For Docker builds, the Docker image needs to copy, build and symlink SwissEph

        # Set up Swisseph libraries
        RUN rm -rf /usr/local/lib/pkgconfig/swisseph.pc && \
            rm -rf /usr/local/include/swisseph && \
            rm -rf /usr/local/lib/swisseph
        # Copy and build Swisseph
        COPY ./ext/swisseph /swisseph
        RUN cd /swisseph/original_project && make clean libswe.a && \
            cp /swisseph/original_project/libswe.a /swisseph/lib/libswe.a && \
            mkdir -p /usr/local/lib/pkgconfig && \
            ln -sf /swisseph/swisseph.pc /usr/local/lib/pkgconfig/swisseph.pc && \
            ln -sf /swisseph/include /usr/local/include/swisseph && \
            ln -sf /swisseph/lib /usr/local/lib/swisseph


## Usage
For example, here is William Lilly's "Considerations Before Judgement" implemented with SacredStar

        package main

        import (
            "flag"
            "github.com/afjoseph/sacredstar/chart"
            "github.com/afjoseph/sacredstar/house"
            "github.com/afjoseph/sacredstar/pointid"
            "github.com/afjoseph/sacredstar/rulership"
            "github.com/afjoseph/sacredstar/timeandzone"
            "github.com/afjoseph/sacredstar/wrapper"
            "github.com/afjoseph/sacredstar/zodiacalpos"
        )

        func main() {
            flag.Parse()

            // See "Setting up SwissEph" in the README
            swissEphPathFlag := flag.String("swisseph-path", "ext/swisseph/", "Path to the SwissEph library")

            // Init wrapper
            mysacredstar := wrapper.NewWithPath(*swissEphPathFlag)
            defer mysacredstar.Close()

            chartTime := time.Date(2025, 1, 9, 10, 21, 0, 0, time.UTC)
            currTimeInJulian := mysacredstar.GoTimeToJulianDay(date)
            lon, lat, _ := timeandzone.Zones.GetLonLatFromTimezone("Europe/Zurich")
            myChart, _ := chart.NewChartFromJulianDay(
                mysacredstar,
                currTimeInJulian,
                lon, lat,
                chart.TropicalChartType,
                // Or pointid.ModernPlanets, or pointid.VedicPlanets
                pointid.TraditionalPlanets,
            )

            // Ascendant too early / too late
            ascendant := myChart.MustGetPoint(pointid.ASC)
            if ascendant.ZodiacalPos.Degrees < 4 {
                ...
            }
            if ascendant.ZodiacalPos.Degrees > 26 {
                ...
            }

            // Moon too late
            moon := myChart.MustGetPoint(pointid.Moon)
            if moon.ZodiacalPos.Degrees > 26 {
                ...
            }

            // Seventh house afflicted
            seventhHouseRuler := rulership.MustFindRuler(myChart, house.House7)
            seventhHouseRulerStrength, s, reasons := rulership.GetStrength(
                myChart,
                seventhHouseRuler,
            )
            if seventhHouseRulerStrength == rulership.StrengthWeak {
                ...
            }

            // Ascendant ruler combust
            ascendantRuler := rulership.MustFindRuler(myChart, house.House1)
            sunPoint := myChart.MustGetPoint(pointid.Sun)
            asp := ascendantRuler.ZodiacalPos.GetAspect(sunPoint.ZodiacalPos)
            if asp != nil && asp.Type == zodiacalpos.AspectType_Conjunction {
                ...
            }

            // Saturn on ascendant or descendant
            // This means Saturn is conjunct/opposed with the ascendant
            saturn := myChart.MustGetPoint(pointid.Saturn)
            asp = saturn.ZodiacalPos.GetAspect(ascendant.ZodiacalPos)
            if asp != nil &&
                (asp.Type == zodiacalpos.AspectType_Conjunction ||
                    asp.Type == zodiacalpos.AspectType_Opposition) {
                ...
            }
        }

## Testing
Run tests with `just test`. This is a really good way to see if your machine's integration is sound.
