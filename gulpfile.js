// gulpと使用するプラグインを読み込む
var gulp = require("gulp"),
    concat = require("gulp-concat"),
    rename = require('gulp-rename'),
    less = require('gulp-less'),
    minifyCss = require("gulp-minify-css"),
    uglify = require("gulp-uglify"),
    gulpFilter = require('gulp-filter'),
    bower = require('main-bower-files'),
    del = require('del');

// Common config
var config = {
    src: './src',
    dist: './dist',
    isWatchify: false,
    isRelease: false
};

// bowerで導入したパッケージのCSSを取ってくるタスク
gulp.task('css', function () {
    var cssLibDir = 'dev/css/lib/', // cssを出力するディレクトリ
        cssFilter = gulpFilter('**/*.css', { restore: true }),
        lessFilter = gulpFilter('**/*.less', { restore: true }); // Bootstrapのコアがlessなのでlessをファイルを抽出するフィルター
    return gulp.src(bower({ paths: { bowerJson: 'bower.json' } }))
        .pipe(cssFilter)
        .pipe(rename({ prefix: '_', extname: '.css' }))
        .pipe(gulp.dest(cssLibDir))
        .pipe(cssFilter.restore)
        .pipe(lessFilter)
        .pipe(less())
        .pipe(rename({ prefix: '_', extname: '.css' }))
        .pipe(gulp.dest(cssLibDir))
        .pipe(lessFilter.restore);
});

// パッケージのCSSを1つに結合してmin化するタスク
gulp.task('css.concat', ['css'], function () {
    var cssDir = 'dev/css/', // 結合したcssを出力するディレクトリ
        cssLibDir = 'dev/css/lib/'; // ライブラリのCSSが置いてあるディレクトリ
    return gulp.src(cssLibDir + '_*.css')
        .pipe(concat('style.css'))
        .pipe(gulp.dest(cssDir))
        .pipe(minifyCss())
        .pipe(rename({ extname: '.min.css' }))
        .pipe(gulp.dest(cssDir));
});

gulp.task('css.copy', ['css.concat'], function () { return gulp.src(["dev/css/style.css"]).pipe(gulp.dest('css/')); });
gulp.task('css.release', function () { return gulp.src(["dev/css/style.min.css"]).pipe(gulp.dest('css/')).pipe(rename('style.css')); });

// bowerで導入したパッケージのjsを取ってきて1つにまとめるタスク
gulp.task('js', function () {
    var jsDir = 'dev/js/', // jsを出力するディレクトリ
        jsFilter = gulpFilter('**/*.js', { restore: true }); // jsファイルを抽出するフィルター
    return gulp.src(bower({
        paths: {
            bowerJson: 'bower.json'
        }
    }))
        .pipe(jsFilter)
        .pipe(concat('script.js'))
        // jsを1つにしたものを出力
        .pipe(gulp.dest(jsDir))
        .pipe(uglify({ output: { comments: /^!/ } }))
        .pipe(rename({ extname: '.min.js' }))
        // jsを1つにしてmin化したものを出力
        .pipe(gulp.dest(jsDir))
        .pipe(jsFilter.restore);
});

gulp.task('js.copy', ['js'], function () { return gulp.src(["dev/js/script.js"]).pipe(gulp.dest('js/')); });
gulp.task('js.release', ['js'], function () { return gulp.src(["dev/js/script.min.js"]).pipe(gulp.dest('js/')).pipe(rename('script.js')); });

gulp.task('Copy.fonts', function () { return gulp.src(["bower_components/bootstrap/dist/fonts/*"]).pipe(gulp.dest('fonts/')); });

gulp.task('build', ['css', 'css.concat', 'css.copy', 'js', 'js.copy', 'Copy.fonts']);
gulp.task('release', ['css', 'css.concat', 'js', 'css.release', 'js.release', 'Copy.fonts']);
gulp.task('clean', del.bind(null, ['fonts/*', 'dev/*', 'bower_components/*', 'fonts', 'dev', 'bower_components', 'js/script.js', 'css/style.css'], { dot: true }));
