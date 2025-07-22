import { assert } from "chai";
import { readFileSync } from "fs";
import path from "path";
import { PemEntry } from "./pemEntry";
import { USER_SEED_LENGTH } from "./userKeys";

describe("test pem entry", () => {
    const walletsPath = path.join("src", "testdata", "testwallets");

    it("should create from text all", () => {
        let text = readFileSync(path.join(walletsPath, "alice.pem"), "utf-8");
        let entries = PemEntry.fromTextAll(text);
        let entry = entries[0];

        assert.lengthOf(entries, 1);
        assert.equal(entry.label, "drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh");
        assert.equal(
            Buffer.from(entry.message.slice(0, USER_SEED_LENGTH)).toString("hex"),
            "2bbcdae7e193924fa0d301e7a12c7defc92a93bc5e587cc968f04fcb86022e1c",
        );

        text = readFileSync(path.join(walletsPath, "multipleUserKeys.pem"), "utf-8");
        entries = PemEntry.fromTextAll(text);
        entry = entries[0];

        assert.lengthOf(entries, 3);
        assert.equal(entry.label, "drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh");
        assert.equal(
            Buffer.from(entry.message.slice(0, USER_SEED_LENGTH)).toString("hex"),
            "2bbcdae7e193924fa0d301e7a12c7defc92a93bc5e587cc968f04fcb86022e1c",
        );

        entry = entries[1];
        assert.equal(entry.label, "drt18h03w0y7qtqwtra3u4f0gu7e3kn2fslj83lqxny39m5c4rwaectswerhd2");
        assert.equal(
            Buffer.from(entry.message.slice(0, USER_SEED_LENGTH)).toString("hex"),
            "8928add00f0d168620a76ec7af31a92f957038a1a2ed75778a4243248d319f2f",
        );

        entry = entries[2];
        assert.equal(entry.label, "drt1kp072dwz0arfz8m5lzmlypgu2nme9l9q33aty0znualvanfvmy5qd3yy8q");
        assert.equal(
            Buffer.from(entry.message.slice(0, USER_SEED_LENGTH)).toString("hex"),
            "5aa2311a2274ff47cc804f12a4e8b28cf74650a0d1efcb8175b90eda4e3e6b4c",
        );
    });

    it("should create from text all for validators", () => {
        let text = readFileSync(path.join(walletsPath, "validatorKey00.pem"), "utf-8");
        let entries = PemEntry.fromTextAll(text);
        let entry = entries[0];

        assert.lengthOf(entries, 1);
        assert.equal(
            entry.label,
            "d3e0427c22ff9cc80ef4156f976644cfa25c54e5a69ed199132053f8cbbfddd4eb15a2f732a3c9b392169c8b1d060e0b5ab0d88b4dd7b4010fa051a17ef81bdbace5e68025965b00bf48e14a9ec8d8e2a8bcc9e62f97ddac3268f6b805f7b80e",
        );
        assert.equal(
            Buffer.from(entry.message.slice(0, USER_SEED_LENGTH)).toString("hex"),
            "132e9b47291fcc62c64b334fd434ab2db74bf64b42d4cc1b4cedd10df77c1936",
        );

        text = readFileSync(path.join(walletsPath, "multipleValidatorKeys.pem"), "utf-8");
        entries = PemEntry.fromTextAll(text);
        entry = entries[0];

        assert.lengthOf(entries, 4);
        assert.equal(
            entry.label,
            "d3e0427c22ff9cc80ef4156f976644cfa25c54e5a69ed199132053f8cbbfddd4eb15a2f732a3c9b392169c8b1d060e0b5ab0d88b4dd7b4010fa051a17ef81bdbace5e68025965b00bf48e14a9ec8d8e2a8bcc9e62f97ddac3268f6b805f7b80e",
        );
        assert.equal(
            Buffer.from(entry.message).toString("hex"),
            "132e9b47291fcc62c64b334fd434ab2db74bf64b42d4cc1b4cedd10df77c1936",
        );

        entry = entries[1];
        assert.equal(
            entry.label,
            "b0b6349b3f693e08c433970d10efb2fe943eac4057a945146bee5fd163687f4e1800d541aa0f11bf9e4cb6552f512e126068e68eb471d18fcc477ddfe0b9b3334f34e30d8b7b2c08f914f4ae54454f75fb28922ba9fd28785bcadc627031fa8a",
        );
        assert.equal(
            Buffer.from(entry.message).toString("hex"),
            "490a55ab40cbea7798f7a743bcd3da44276bec6ad083759a4516646186257466",
        );

        entry = entries[2];
        assert.equal(
            entry.label,
            "67c301358a41bef74df2ae6aa9914e3a5e7a4b528bbd19596cca4b2fd97a62ab2c0a88b02adf1c5973a82c7544cdc40539ae62a9ac05351cfc59c300bbf4492f4266c550987355c39cff8e84ff74e012c7fd372c240eeb916ef87eead82ffd98",
        );
        assert.equal(
            Buffer.from(entry.message).toString("hex"),
            "59074392f4f970c3b5d4ba1783955fcf9fcc44d9a5aff2b5cddb0020ca155b5c",
        );

        entry = entries[3];
        assert.equal(
            entry.label,
            "ab0a22ba2be6560af8520208393381760f9d4f69fca4f152b0a3fe7b124dd7f932fd8c1fbb372792c235baafac36030ceaf6ebf215de4e8d8d239f347f2fed10a75a07cbf9dc56efbbfca2e319152a363df122c300cdeb2faa02a61ebefd8a0e",
        );
        assert.equal(
            Buffer.from(entry.message).toString("hex"),
            "a6091fb517645196396d04aab3679e4f059db9d85c981625c1999dada985cd5e",
        );
    });

    it("should create from text all for validators with extra lines in pem file", () => {
        let text = readFileSync(path.join(walletsPath, "validatorKey00WithExtraLines.pem"), "utf-8");
        let entries = PemEntry.fromTextAll(text);
        let entry = entries[0];

        assert.lengthOf(entries, 1);
        assert.equal(
            entry.label,
            "d3e0427c22ff9cc80ef4156f976644cfa25c54e5a69ed199132053f8cbbfddd4eb15a2f732a3c9b392169c8b1d060e0b5ab0d88b4dd7b4010fa051a17ef81bdbace5e68025965b00bf48e14a9ec8d8e2a8bcc9e62f97ddac3268f6b805f7b80e",
        );
        assert.equal(
            Buffer.from(entry.message.slice(0, USER_SEED_LENGTH)).toString("hex"),
            "132e9b47291fcc62c64b334fd434ab2db74bf64b42d4cc1b4cedd10df77c1936",
        );

        text = readFileSync(path.join(walletsPath, "multipleValidatorKeys.pem"), "utf-8");
        entries = PemEntry.fromTextAll(text);
        entry = entries[0];

        assert.lengthOf(entries, 4);
        assert.equal(
            entry.label,
            "d3e0427c22ff9cc80ef4156f976644cfa25c54e5a69ed199132053f8cbbfddd4eb15a2f732a3c9b392169c8b1d060e0b5ab0d88b4dd7b4010fa051a17ef81bdbace5e68025965b00bf48e14a9ec8d8e2a8bcc9e62f97ddac3268f6b805f7b80e",
        );
        assert.equal(
            Buffer.from(entry.message).toString("hex"),
            "132e9b47291fcc62c64b334fd434ab2db74bf64b42d4cc1b4cedd10df77c1936",
        );

        entry = entries[1];
        assert.equal(
            entry.label,
            "b0b6349b3f693e08c433970d10efb2fe943eac4057a945146bee5fd163687f4e1800d541aa0f11bf9e4cb6552f512e126068e68eb471d18fcc477ddfe0b9b3334f34e30d8b7b2c08f914f4ae54454f75fb28922ba9fd28785bcadc627031fa8a",
        );
        assert.equal(
            Buffer.from(entry.message).toString("hex"),
            "490a55ab40cbea7798f7a743bcd3da44276bec6ad083759a4516646186257466",
        );

        entry = entries[2];
        assert.equal(
            entry.label,
            "67c301358a41bef74df2ae6aa9914e3a5e7a4b528bbd19596cca4b2fd97a62ab2c0a88b02adf1c5973a82c7544cdc40539ae62a9ac05351cfc59c300bbf4492f4266c550987355c39cff8e84ff74e012c7fd372c240eeb916ef87eead82ffd98",
        );
        assert.equal(
            Buffer.from(entry.message).toString("hex"),
            "59074392f4f970c3b5d4ba1783955fcf9fcc44d9a5aff2b5cddb0020ca155b5c",
        );

        entry = entries[3];
        assert.equal(
            entry.label,
            "ab0a22ba2be6560af8520208393381760f9d4f69fca4f152b0a3fe7b124dd7f932fd8c1fbb372792c235baafac36030ceaf6ebf215de4e8d8d239f347f2fed10a75a07cbf9dc56efbbfca2e319152a363df122c300cdeb2faa02a61ebefd8a0e",
        );
        assert.equal(
            Buffer.from(entry.message).toString("hex"),
            "a6091fb517645196396d04aab3679e4f059db9d85c981625c1999dada985cd5e",
        );
    });
    it("should convert to text", () => {
        let text = readFileSync(path.join(walletsPath, "alice.pem"), "utf-8").trim();
        assert.deepEqual(PemEntry.fromTextAll(text)[0].toText(), text);

        text = readFileSync(path.join(walletsPath, "validatorKey00.pem"), "utf-8").trim();
        assert.deepEqual(PemEntry.fromTextAll(text)[0].toText(), text);

        text = readFileSync(path.join(walletsPath, "multipleUserKeys.pem"), "utf-8").trim();
        let entries = PemEntry.fromTextAll(text);
        let actualText = entries.map((entry) => entry.toText()).join("\n");
        assert.deepEqual(actualText, text);

        text = readFileSync(path.join(walletsPath, "multipleValidatorKeys.pem"), "utf-8").trim();
        entries = PemEntry.fromTextAll(text);
        actualText = entries.map((entry) => entry.toText()).join("\n");
        assert.deepEqual(actualText, text);
    });
});
