package makemkvdriver

import (
	"strings"
	"testing"
)

func TestXCopyParser(t *testing.T) {
	t.Log("Begin XCopyParser Test")
	const input = `

-------------------------------------------------------------------------------
   ROBOCOPY     ::     Robust File Copy for Windows                              
-------------------------------------------------------------------------------

  Started : Saturday, January 9, 2016 9:04:08 PM
   Source : C:\Users\james\Desktop\iPlayerRecordings\
     Dest : C:\Users\james\Desktop\iPlayerRecordings2\

    Files : *.*
	    
  Options : *.* /DCOPY:DA /COPY:DAT /R:1000000 /W:30 

------------------------------------------------------------------------------

	  New Dir          1	C:\Users\james\Desktop\iPlayerRecordings\
	    New File  		  37.7 m	A_Festival_of_Nine_Lessons_and_Carols_-_25_12_2015_b06s9hzg_default.m4a  0.0%  2.6%  5.2%  7.9% 10.5% 13.2% 15.8% 18.5% 21.1% 23.8% 26.4% 29.1% 31.7% 34.4% 37.0% 39.7% 42.3% 45.0% 47.6% 50.3% 52.9% 55.6% 58.2% 60.9% 63.5% 66.2% 68.8% 71.5% 74.1% 76.8% 79.4% 82.0% 84.7% 87.3% 90.0% 92.6% 95.3% 97.9%100%  

------------------------------------------------------------------------------

               Total    Copied   Skipped  Mismatch    FAILED    Extras
    Dirs :         1         1         0         0         0         0
   Files :         1         1         0         0         0         0
   Bytes :   37.76 m   37.76 m         0         0         0         0
   Times :   0:00:00   0:00:00                       0:00:00   0:00:00


   Speed :           1164557882 Bytes/sec.
   Speed :            66636.536 MegaBytes/min.
   Ended : Saturday, January 9, 2016 9:04:08 PM
	`
	stream := strings.NewReader(input)
	dc := roboCommand{}

	reader := func(value int) {
		t.Log(value)
	}

	status := func(val string) {
		t.Log("Line:" + val)
	}
	dc.parseStdOut(stream, reader, status)
}
