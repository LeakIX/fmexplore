## FMExplore

FMExplore ( FileMaker explore ) is a tool for exploring and dumping Filemaker instances.

### Usage

```bash
$ ./fmexplore  https://username:password@xyz.fmcloud.fm ./fm-dump-directory
2021/09/15 01:35:12 Found database PersonalityTest
2021/09/15 01:35:12 Dumping database PersonalityTest, layout PersonalityTest into fm-dump-directory/PersonalityTest-PersonalityTest.json
2021/09/15 01:35:12 Dumping database PersonalityTest, layout web into fm-dump-directory/PersonalityTest-web.json
2021/09/15 01:35:12 Dumping database PersonalityTest, layout Report Page1 into fm-dump-directory/PersonalityTest-Report Page1.json
2021/09/15 01:35:12 Dumping database PersonalityTest, layout Report Page2 into fm-dump-directory/PersonalityTest-Report Page2.json
2021/09/15 01:35:12 Dumping database PersonalityTest, layout Report Page3 into fm-dump-directory/PersonalityTest-Report Page3.json
2021/09/15 01:35:13 Dumping database PersonalityTest, layout Report Page4 into fm-dump-directory/PersonalityTest-Report Page4.json
2021/09/15 01:35:13 Dumping database PersonalityTest, layout Report Page5 into fm-dump-directory/PersonalityTest-Report Page5.json
2021/09/15 01:35:13 Dumping database PersonalityTest, layout Report Page6 into fm-dump-directory/PersonalityTest-Report Page6.json
2021/09/15 01:35:13 Dumping database PersonalityTest, layout Settings into fm-dump-directory/PersonalityTest-Settings.json
```

### Output

Every database and layout is saved in a separate `json` file in the output directory.

If a field contains JSON data, it's parsed as well.