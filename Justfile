swisseph_path := justfile_directory() + "/ext/swisseph/"

test:
    go test -race -failfast -count=1 -v ./...

build-and-symlink-swisseph:
    rm -rf /usr/local/lib/pkgconfig/swisseph.pc
    rm -rf /usr/local/include/swisseph
    rm -rf /usr/local/lib/swisseph
    (cd {{ swisseph_path }}/original_project && make clean libswe.a)
    cp {{ swisseph_path }}/original_project/libswe.a {{ swisseph_path }}/lib/libswe.a
    ln -sf {{ swisseph_path }}/swisseph.pc /usr/local/lib/pkgconfig/swisseph.pc
    ln -sf {{ swisseph_path }}/include /usr/local/include/swisseph
    ln -sf {{ swisseph_path }}/lib /usr/local/lib/swisseph
