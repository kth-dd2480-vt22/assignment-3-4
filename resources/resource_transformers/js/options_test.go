// Copyright 2020 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package js

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/hugofs"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/media"

	"github.com/evanw/esbuild/pkg/api"

	qt "github.com/frankban/quicktest"
)

// This test is added to test/warn against breaking the "stability" of the
// cache key. It's sometimes needed to break this, but should be avoided if possible.
func TestOptionKey(t *testing.T) {
	c := qt.New(t)

	opts := map[string]interface{}{
		"TargetPath": "foo",
		"Target":     "es2018",
	}

	key := (&buildTransformation{optsm: opts}).Key()

	c.Assert(key.Value(), qt.Equals, "jsbuild_7891849149754191852")
}

func TestToBuildOptions(t *testing.T) {
	fmt.Println("TestToBuildOptions")
	c := qt.New(t)

	var cover [100]int
	var temp [100]int

	// Test 1
	fmt.Println("Test 1")
	/* 
	 * This one tests when mediaType is set to media.JavascriptType,
	 * and the rest of the options are set to default values 
	 * -> branch 0, 10, 11, 16, 21 and 24
	 */
	opts, err, temp := toBuildOptions(Options{mediaType: media.JavascriptType})
	for i := 0; i < 100; i++ {
		cover[i] += temp[i]
	}

	c.Assert(err, qt.IsNil)
	c.Assert(opts, qt.DeepEquals, api.BuildOptions{
		Bundle: true,
		Target: api.ESNext,
		Format: api.FormatIIFE,
		Stdin: &api.StdinOptions{
			Loader: api.LoaderJS,
		},
	})

	for i := 0; i < 26; i++ {
		if temp[i] > 0 {
			fmt.Printf("branch %d is covered %d times\n", i, temp[i])
		}
	}

	// Test 2
	fmt.Println("Test 2")
	/*
	 * This one tests when Target is set to "es2018", Format is set to "cjs",
	 * Minify is set to true, mediaType is set to media.JavascriptType,
	 * , AvoidTDZ is set to true, and the rest set to default values
	 * -> branch 5, 10, 11, 18, 21 and 24
	 */
	opts, err, temp = toBuildOptions(Options{
		Target:    "es2018",
		Format:    "cjs",
		Minify:    true,
		mediaType: media.JavascriptType,
		AvoidTDZ:  true,
	})
	c.Assert(err, qt.IsNil)
	c.Assert(opts, qt.DeepEquals, api.BuildOptions{
		Bundle:            true,
		Target:            api.ES2018,
		Format:            api.FormatCommonJS,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		MinifyWhitespace:  true,
		Stdin: &api.StdinOptions{
			Loader: api.LoaderJS,
		},
	})
	for i := 0; i < 100; i++ {
		cover[i] += temp[i]
	}
	for i := 0; i < 26; i++ {
		if temp[i] > 0 {
			fmt.Printf("branch %d is covered %d times\n", i, temp[i])
		}
	}

	// Test 3
	fmt.Println("Test 3")
	/*
	 * This one tests when Target is set to "es2018", Format is set to "cjs",
	 * Minify is set to true, mediaType is set to media.JavascriptType,
	 * , SourceMap set to "inline", and the rest set to default values
	 * -> branch 5, 10, 11, 18, 21 and 22
	 */
	opts, err, temp = toBuildOptions(Options{
		Target: "es2018", Format: "cjs", Minify: true, mediaType: media.JavascriptType,
		SourceMap: "inline",
	})
	c.Assert(err, qt.IsNil)
	c.Assert(opts, qt.DeepEquals, api.BuildOptions{
		Bundle:            true,
		Target:            api.ES2018,
		Format:            api.FormatCommonJS,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		MinifyWhitespace:  true,
		Sourcemap:         api.SourceMapInline,
		Stdin: &api.StdinOptions{
			Loader: api.LoaderJS,
		},
	})
	for i := 0; i < 100; i++ {
		cover[i] += temp[i]
	}
	for i := 0; i < 26; i++ {
		if temp[i] > 0 {
			fmt.Printf("branch %d is covered %d times\n", i, temp[i])
		}
	}

	// Test 4
	fmt.Println("Test 4")
	/*
	 * This one tests when Target is set to "es2018", Format is set to "cjs",
	 * Minify is set to true, mediaType is set to media.JavascriptType,
	 * , SourceMap set to "inline", and the rest set to default values
	 * -> branch 5, 10, 11, 18, 21 and 22
	 */
	opts, err, temp = toBuildOptions(Options{
		Target: "es2018", Format: "cjs", Minify: true, mediaType: media.JavascriptType,
		SourceMap: "inline",
	})
	c.Assert(err, qt.IsNil)
	c.Assert(opts, qt.DeepEquals, api.BuildOptions{
		Bundle:            true,
		Target:            api.ES2018,
		Format:            api.FormatCommonJS,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		MinifyWhitespace:  true,
		Sourcemap:         api.SourceMapInline,
		Stdin: &api.StdinOptions{
			Loader: api.LoaderJS,
		},
	})
	for i := 0; i < 100; i++ {
		cover[i] += temp[i]
	}
	for i := 0; i < 26; i++ {
		if temp[i] > 0 {
			fmt.Printf("branch %d is covered %d times\n", i, temp[i])
		}
	}

	// Test 5
	fmt.Println("Test 5")
	/*
	 * This one tests when Target is set to "es2018", Format is set to "cjs",
	 * Minify is set to true, mediaType is set to media.JavascriptType,
	 * , SourceMap set to "external", and the rest set to default values
	 * -> branch 5, 10, 11, 18, 21 and 23
	 */
	opts, err, temp = toBuildOptions(Options{
		Target: "es2018", Format: "cjs", Minify: true, mediaType: media.JavascriptType,
		SourceMap: "external",
	})
	c.Assert(err, qt.IsNil)
	c.Assert(opts, qt.DeepEquals, api.BuildOptions{
		Bundle:            true,
		Target:            api.ES2018,
		Format:            api.FormatCommonJS,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		MinifyWhitespace:  true,
		Sourcemap:         api.SourceMapExternal,
		Stdin: &api.StdinOptions{
			Loader: api.LoaderJS,
		},
	})
	for i := 0; i < 100; i++ {
		cover[i] += temp[i]
	}
	for i := 0; i < 26; i++ {
		if temp[i] > 0 {
			fmt.Printf("branch %d is covered %d times\n", i, temp[i])
		}
	}
	
	// Test 6
	fmt.Println("Test 6")
	/*
	 * This one tests when mediaType is set to media.JavascriptType,
	 * , SourceMap set to "error", and the rest set to default values
	 * -> branch 0, 10, 11, 16, 21, 25
	 */
	opts, err, temp = toBuildOptions(Options{
		mediaType: media.JavascriptType, SourceMap: "error",
	})
	c.Assert(err, qt.IsNotNil)
	for i := 0; i < 100; i++ {
		cover[i] += temp[i]
	}
	for i := 0; i < 26; i++ {
		if temp[i] > 0 {
			fmt.Printf("branch %d is covered %d times\n", i, temp[i])
		}
	}

	// Test 7
	fmt.Println("Test 7")
	/*
	 * This one tests when Target is "error", 
	 * and the rest set to default values
	 * -> branch 8
	 */
	opts, err, temp = toBuildOptions(Options{
		Target: "error",
	})
	c.Assert(err, qt.IsNotNil)
	for i := 0; i < 100; i++ {
		cover[i] += temp[i]
	}
	for i := 0; i < 26; i++ {
		if temp[i] > 0 {
			fmt.Printf("branch %d is covered %d times\n", i, temp[i])
		}
	}

	// Test 8
	fmt.Println("Test 8")
	/*
	 * This one tests when mediaType is set to media.HTMLType,
	 * and the rest set to default values
	 * -> branch 0, 10, 15
	 */
	opts, err, temp = toBuildOptions(Options{
		mediaType: media.HTMLType,
	})
	c.Assert(err, qt.IsNotNil)
	for i := 0; i < 100; i++ {
		cover[i] += temp[i]
	}
	for i := 0; i < 26; i++ {
		if temp[i] > 0 {
			fmt.Printf("branch %d is covered %d times\n", i, temp[i])
		}
	}

	// Test 9
	fmt.Println("Test 9")
	/*
	 * This one tests when Format is "error",
	 * and the rest set to default values
	 * -> branch 0, 9, 11, 19
	 */
	opts, err, temp = toBuildOptions(Options{
		Format: "error",
	})
	c.Assert(err, qt.IsNotNil)
	for i := 0; i < 100; i++ {
		cover[i] += temp[i]
	}
	for i := 0; i < 26; i++ {
		if temp[i] > 0 {
			fmt.Printf("branch %d is covered %d times\n", i, temp[i])
		}
	}

	// Test 10
	fmt.Println("Test 10")
	/*
	 * This one tests when Target is set to "es5", mediaType is TypeScriptType,
	 * Format is set to "esm", Defines is map[string]interface{}{"one": 1,}, 
	 * and the rest set to default values
	 * -> branch 1, 10, 12, 17, 20, 24
	 */
	opts, err, temp = toBuildOptions(Options{
		Target: "es5", mediaType: media.TypeScriptType, Format: "esm",
		Defines: map[string]interface{}{"one": 1,},
	})
	c.Assert(err, qt.IsNil)
	c.Assert(opts, qt.DeepEquals, api.BuildOptions{
		Bundle:            true,
		Target:            api.ES5,
		Format:            api.FormatESModule,
		Define: 		   maps.ToStringMapString(map[string]interface{}{"one": 1,}),
		Stdin: &api.StdinOptions{
			Loader: api.LoaderTS,
		},
	})
	
	for i := 0; i < 100; i++ {
		cover[i] += temp[i]
	}
	for i := 0; i < 26; i++ {
		if temp[i] > 0 {
			fmt.Printf("branch %d is covered %d times\n", i, temp[i])
		}
	}

	// Test 11
	fmt.Println("Test 11")
	/*
	 * This one tests when Target is set to "es6", mediaType is TSXType,
	 * and the rest set to default values
	 * -> branch 2, 10, 13, 16, 21, 24
	 */
	opts, err, temp = toBuildOptions(Options{
		Target: "es6", mediaType: media.TSXType,})
	c.Assert(err, qt.IsNil)
	c.Assert(opts, qt.DeepEquals, api.BuildOptions{
		Bundle:            true,
		Target:            api.ES2015,
		Format:			   api.FormatIIFE,
		Stdin: &api.StdinOptions{
			Loader: api.LoaderTSX,
		},
	})
	
	for i := 0; i < 100; i++ {
		cover[i] += temp[i]
	}
	for i := 0; i < 26; i++ {
		if temp[i] > 0 {
			fmt.Printf("branch %d is covered %d times\n", i, temp[i])
		}
	}

	// Test 12
	fmt.Println("Test 12")
	/*
	 * This one tests when Target is set to "es2016", mediaType is JSXType,
	 * and the rest set to default values
	 * -> branch 2, 10, 13, 16, 21, 24
	 */
	opts, err, temp = toBuildOptions(Options{
		Target: "es2016", mediaType: media.JSXType,})
	c.Assert(err, qt.IsNil)
	c.Assert(opts, qt.DeepEquals, api.BuildOptions{
		Bundle:            true,
		Target:            api.ES2016,
		Format:			   api.FormatIIFE,
		Stdin: &api.StdinOptions{
			Loader: api.LoaderJSX,
		},
	})
	for i := 0; i < 100; i++ {
		cover[i] += temp[i]
	}
	for i := 0; i < 26; i++ {
		if temp[i] > 0 {
			fmt.Printf("branch %d is covered %d times\n", i, temp[i])
		}
	}

	// Test 13
	fmt.Println("Test 13")
	/*
	 * This one tests when Target is set to "es2016", mediaType is JSXType,
	 * and the rest set to default values
	 * -> branch 4, 9, 11, 16, 21, 24
	 */
	opts, err, temp = toBuildOptions(Options{
		Target: "es2017", })
	c.Assert(err, qt.IsNil)
	c.Assert(opts, qt.DeepEquals, api.BuildOptions{
		Bundle:            true,
		Target:            api.ES2017,
		Format:			   api.FormatIIFE,
		Stdin: &api.StdinOptions{
			Loader: api.LoaderJS,
		},
	})
	for i := 0; i < 100; i++ {
		cover[i] += temp[i]
	}
	for i := 0; i < 26; i++ {
		if temp[i] > 0 {
			fmt.Printf("branch %d is covered %d times\n", i, temp[i])
		}
	}

	// Test 14
	fmt.Println("Test 14")
	/*
		* This one tests when Target is set to "es2016", mediaType is JSXType,
		* and the rest set to default values
		* -> branch 6, 9, 11, 16, 21, 24
		*/
	opts, err, temp = toBuildOptions(Options{
		Target: "es2019", })
	c.Assert(err, qt.IsNil)
	c.Assert(opts, qt.DeepEquals, api.BuildOptions{
		Bundle:            true,
		Target:            api.ES2019,
		Format:			   api.FormatIIFE,
		Stdin: &api.StdinOptions{
			Loader: api.LoaderJS,
		},
	})
	for i := 0; i < 100; i++ {
		cover[i] += temp[i]
	}
	for i := 0; i < 26; i++ {
		if temp[i] > 0 {
			fmt.Printf("branch %d is covered %d times\n", i, temp[i])
		}
	}

	// Test 15
	fmt.Println("Test 15")
	/*
		* This one tests when Target is set to "es2016", mediaType is JSXType,
		* and the rest set to default values
		* -> branch 7, 9, 11, 16, 21, 24
		*/
	opts, err, temp = toBuildOptions(Options{
		Target: "es2020", })
	c.Assert(err, qt.IsNil)
	c.Assert(opts, qt.DeepEquals, api.BuildOptions{
		Bundle:            true,
		Target:            api.ES2020,
		Format:			   api.FormatIIFE,
		Stdin: &api.StdinOptions{
			Loader: api.LoaderJS,
		},
	})
	for i := 0; i < 100; i++ {
		cover[i] += temp[i]
	}
	for i := 0; i < 26; i++ {
		if temp[i] > 0 {
			fmt.Printf("branch %d is covered %d times\n", i, temp[i])
		}
	}

	// Print the coverage for the 26 branches.
	fmt.Println("Branch Coverage for TestToBuildOptions:")
	var covered int
	for i := 0; i < 26; i++ {
		if cover[i] > 0 {
			covered++
			fmt.Printf("branch %d is covered %d times\n", i, cover[i])
		}
	}
	// Print the percentage of branches covered.
	fmt.Printf("%d branches, which is %.2f%% of the branches are covered\n", covered, float64(covered)/float64(26)*100)

}

func TestResolveComponentInAssets(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		name    string
		files   []string
		impPath string
		expect  string
	}{
		{"Basic, extension", []string{"foo.js", "bar.js"}, "foo.js", "foo.js"},
		{"Basic, no extension", []string{"foo.js", "bar.js"}, "foo", "foo.js"},
		{"Basic, no extension, typescript", []string{"foo.ts", "bar.js"}, "foo", "foo.ts"},
		{"Not found", []string{"foo.js", "bar.js"}, "moo.js", ""},
		{"Not found, double js extension", []string{"foo.js.js", "bar.js"}, "foo.js", ""},
		{"Index file, folder only", []string{"foo/index.js", "bar.js"}, "foo", "foo/index.js"},
		{"Index file, folder and index", []string{"foo/index.js", "bar.js"}, "foo/index", "foo/index.js"},
		{"Index file, folder and index and suffix", []string{"foo/index.js", "bar.js"}, "foo/index.js", "foo/index.js"},

		// Issue #8949
		{"Check file before directory", []string{"foo.js", "foo/index.js"}, "foo", "foo.js"},
	} {

		c.Run(test.name, func(c *qt.C) {
			baseDir := "assets"
			mfs := afero.NewMemMapFs()

			for _, filename := range test.files {
				c.Assert(afero.WriteFile(mfs, filepath.Join(baseDir, filename), []byte("let foo='bar';"), 0777), qt.IsNil)
			}

			bfs := hugofs.DecorateBasePathFs(afero.NewBasePathFs(mfs, baseDir).(*afero.BasePathFs))

			got := resolveComponentInAssets(bfs, test.impPath)

			gotPath := ""
			if got != nil {
				gotPath = filepath.ToSlash(got.Path)
			}

			c.Assert(gotPath, qt.Equals, test.expect)
		})

	}
}
