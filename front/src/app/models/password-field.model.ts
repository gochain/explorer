export class PasswordField {
  public type: string;
  public icon: string;
  public alt: string;

  private _isShown: boolean;
  private set isShown(value: boolean) {
    this._isShown = value;
    this.process();
  }

  constructor(isShown = false) {
    this.isShown = isShown;
  }

  public toggle(): void {
    this.isShown = !this._isShown;
  }

  private process(): void {
    if (this._isShown) {
      this.type = 'text';
      this.icon = 'eye.svg';
      this.alt = 'hide';
    } else {
      this.type = 'password';
      this.icon = 'eye-off.svg';
      this.alt = 'show';
    }
  }
}
