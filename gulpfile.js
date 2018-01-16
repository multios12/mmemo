// gulpと使用するプラグインを読み込む
var gulp = require('gulp');
var tar = require("gulp-tar");
var gzip = require("gulp-gzip");
var del = require('del');

gulp.task('releaseDocker', function () {
    gulp.src(['bin/www', 'app.js','routes/*', 'dist/*','dist/views/*', 'package.json', 'package-lock.json'], { base: './' })
        .pipe(tar('hmemo.tar'))
        .pipe(gzip())
        .pipe(gulp.dest('./docker'));
});

gulp.task('clean', del.bind(null, [ './docker/*.tar.gz'], { dot: true }));
