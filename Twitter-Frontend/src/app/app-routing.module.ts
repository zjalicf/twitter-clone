import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ChangePasswordComponent } from './components/change-password/change-password.component';
import { LoginComponent } from './components/login/login.component';
import { MainPageComponent } from './components/main-page/main-page.component';
import { RecoveryEnterMailComponent } from './components/recovery-enter-mail/recovery-enter-mail.component';
import { RecoveryEnterTokenComponent } from './components/recovery-enter-token/recovery-enter-token.component';
import { MyProfileComponent } from './components/my-profile/my-profile.component';
import { RegisterBusinessComponent } from './components/register-business/register-business.component';
import { RegisterRegularComponent } from './components/register-regular/register-regular.component';
import { TestAuthPageComponent } from './components/test-auth-page/test-auth-page.component';
import { VerifyAccountComponent } from './components/verify-account/verify-account.component';
import { RecoveryNewPasswordsComponent } from './components/recovery-new-passwords/recovery-new-passwords.component';
import { TweetAddComponent } from './components/tweet/tweet-add/tweet-add.component';

const routes: Routes = [
  {
    path: 'Main-Page',
    component: MainPageComponent
  },
  {
    path: 'Register-Regular',
    component: RegisterRegularComponent
  },
  {
    path: 'Register-Business',
    component: RegisterBusinessComponent
  },
  {
    path: 'Login',
    component: LoginComponent
  },
  {
    path: 'Verify-Account',
    component: VerifyAccountComponent
  },
  {
    path: 'Request-Recovery',
    component: RecoveryEnterMailComponent
  },
  {
    path: 'Recovery-Token',
    component: RecoveryEnterTokenComponent
  },
  {
    path: 'Recovery-Password',
    component: RecoveryNewPasswordsComponent
  },
  {
    path: 'Auth-Test',
    component: TestAuthPageComponent
  },
  {
    path: 'My-Profile',
    component: MyProfileComponent
  },
  {
    path: 'Change-Password',
    component: ChangePasswordComponent
  },
  {
    path: 'New-Tweet',
    component: TweetAddComponent
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
