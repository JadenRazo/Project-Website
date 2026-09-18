import { afterEach, test } from 'node:test';
import assert from 'node:assert/strict';
import { safeStorage } from '../src/utils/safeStorage.ts';

afterEach(()=>{ delete globalThis.window; });
function nativeStorage(values = {}) {
  const valuesByKey = new Map(Object.entries(values));
  return {
    get length(){return valuesByKey.size;},
    key:index=>[...valuesByKey.keys()][index] ?? null,
    getItem:key=>valuesByKey.get(key) ?? null,
    setItem:(key,value)=>valuesByKey.set(key,value),
    removeItem:key=>valuesByKey.delete(key),
    clear:()=>valuesByKey.clear(),
  };
}
test('opaque-origin storage getters fall back without accessing a parent window',()=>{
  globalThis.window = Object.defineProperty({},'localStorage',{get(){throw new Error('SecurityError');}});
  const storage=safeStorage('localStorage');
  storage.setItem('theme','dark'); assert.equal(storage.getItem('theme'),'dark');
  storage.removeItem('theme'); assert.equal(storage.getItem('theme'),null);
});
for(const operation of ['setItem','removeItem','clear']) {
  test(`failed ${operation} switches all subsequent operations to memory`,()=>{
    const native=nativeStorage({theme:'light',other:'retained'});
    globalThis.window={localStorage:native};
    const storage=safeStorage('localStorage');
    assert.equal(storage.getItem('theme'),'light');
    assert.equal(storage.getItem('other'),'retained');
    native[operation]=()=>{throw new Error('Quota or privacy restriction');};
    if(operation==='setItem') storage.setItem('theme','dark');
    else if(operation==='removeItem') storage.removeItem('theme');
    else storage.clear();
    assert.equal(storage.getItem('theme'),operation==='setItem'?'dark':null);
    assert.equal(storage.getItem('other'),operation==='clear'?null:'retained');
    storage.setItem('next','works'); assert.equal(storage.getItem('next'),'works');
    storage.removeItem('next'); assert.equal(storage.getItem('next'),null);
  });
}
test('successful persistent reads respect deletion and indexed keys',()=>{
  const native=nativeStorage({theme:'dark'});globalThis.window={localStorage:native};
  const storage=safeStorage('localStorage');
  assert.equal(storage.getItem('theme'),'dark');
  native.clear();assert.equal(storage.getItem('theme'),null);
  assert.equal(storage.length,0);assert.equal(storage.key(0),null);
});
