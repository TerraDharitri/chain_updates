import { UsernameUtils } from "src/endpoints/usernames/username.utils";

describe('Username Utils', () => {
  describe('normalizeUsername', () => {
    it('Username correctly normalized', () => {
      expect(UsernameUtils.normalizeUsername('alice')).toStrictEqual('alice.numbat');
      expect(UsernameUtils.normalizeUsername('alice@google.com')).toStrictEqual('alicegooglecom.numbat');
    });
  });
  describe('encodeUserName', () => {
    it('Encode UserName to Hex', () => {
      expect(UsernameUtils.encodeUsername('alice')).toStrictEqual('616c6963652e6e756d626174');
      expect(UsernameUtils.encodeUsername('alice@google.com')).toStrictEqual('616c696365676f6f676c65636f6d2e6e756d626174');
    });
  });
  describe('getContractAddress', () => {
    it('Get Contract Address from normalized userName', () => {
      expect(UsernameUtils.getContractAddress('alice')).toStrictEqual('drt1qqqqqqqqqqqqqpgq2leexk6fwaxlxggzrnkxzruwsjzfcq2mqzgq7hwwcs');
      expect(UsernameUtils.getContractAddress('alice@google.com')).toStrictEqual('drt1qqqqqqqqqqqqqpgqf4pw79l5s9xkslyf5p06egwcnjul95ksqzxqvw9760');
    });
  });
  describe('extractUsernameFromRawBase64', () => {
    it('should return an empty string for invalid inputs', () => {
      expect(UsernameUtils.extractUsernameFromRawBase64("")).toStrictEqual("");
      expect(UsernameUtils.extractUsernameFromRawBase64("invalid base64")).toStrictEqual("");
    });
    it('should work', () => {
      expect(UsernameUtils.extractUsernameFromRawBase64("YWxpY2U=")).toStrictEqual("alice");
      expect(UsernameUtils.extractUsernameFromRawBase64("YWxpY2UubnVtYmF0")).toStrictEqual("alice.numbat");
      expect(UsernameUtils.extractUsernameFromRawBase64("YWxpY2UyLnN1ZmZpeA==")).toStrictEqual("alice2.suffix");
      expect(UsernameUtils.extractUsernameFromRawBase64("dGVzdC4=")).toStrictEqual("test.");
      expect(UsernameUtils.extractUsernameFromRawBase64("YWxpY2UuZHJ0")).toStrictEqual("alice.drt");
    });
  });
});
