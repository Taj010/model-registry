import { genRandomChars } from 'mod-arch-shared';
import {
  translateDisplayNameForK8s,
  checkValidK8sName,
} from '~/concepts/k8s/K8sNameDescriptionField/utils';

jest.mock('mod-arch-shared', () => ({
  ...jest.requireActual<typeof import('mod-arch-shared')>('mod-arch-shared'),
  genRandomChars: jest.fn(() => 'mockrandom'),
}));

const mockedGenRandomChars = genRandomChars as jest.MockedFunction<typeof genRandomChars>;

describe('K8sNameDescriptionField utils', () => {
  beforeEach(() => {
    mockedGenRandomChars.mockReturnValue('mockrandom');
  });

  describe('translateDisplayNameForK8s', () => {
    it('collapses consecutive dashes and strips leading/trailing dashes', () => {
      expect(translateDisplayNameForK8s('My -- Model (v2)')).toBe('my-model-v2');
    });

    it('returns gen-prefixed name when display normalizes to empty', () => {
      expect(translateDisplayNameForK8s('----')).toBe('gen-mockrandom');
      expect(translateDisplayNameForK8s('!!!')).toBe('gen-mockrandom');
    });

    it('does not generate a name when input is empty or whitespace only', () => {
      expect(translateDisplayNameForK8s('')).toBe('');
      expect(translateDisplayNameForK8s('   ')).toBe('');
    });

    it('prepends safePrefix for non-empty translation', () => {
      expect(translateDisplayNameForK8s('My Model', 'usr-')).toBe('usr-my-model');
    });

    it('uses safePrefix plus random suffix when translation is empty', () => {
      expect(translateDisplayNameForK8s('----', 'pre-')).toBe('pre-mockrandom');
    });
  });

  describe('checkValidK8sName', () => {
    it('accepts valid DNS subdomain-style names', () => {
      expect(checkValidK8sName('my-model-v2').valid).toBe(true);
      expect(checkValidK8sName('gen-mockrandom').valid).toBe(true);
    });
  });
});
